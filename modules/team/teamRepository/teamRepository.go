package teamRepository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NatthawutSK/NoTeams-Backend/modules/team"
	"github.com/jmoiron/sqlx"
)

type ITeamRepository interface {
	CreateTeam(userId string, req *team.CreateTeamReq) (*team.CreateTeamRes, error)
	GetTeamById(teamId string) (*team.GetTeamByIdRes, error)
	JoinTeam(req *team.JoinTeamReq) (*team.JoinTeamRes, error)
	GetTeamByUserId(userId string) ([]*team.GetTeamByUserIdRes, error)
	InviteMember(team_id string, req *team.InviteMemberReq) error
	GetMemberTeam(teamId string) ([]*team.GetMemberTeamRes, error)
	DeleteMember(memberId string) error
	GetAboutTeam(teamId string) (*team.GetAboutTeamRes, error)
	GetSettingTeam(teamId string) (*team.GetSettingTeamRes, error)
	UpdateTeam(teamId string, req *team.UpdateTeamReq) error
	UpdatePermission(teamId string, req *team.UpdatePermissionReq) error
	UpdateCodeTeam(teamId string, req *team.UpdateCodeTeamReq) error
	DeleteTeam(teamId string) error
	ExitTeam(userId, teamId string) error
}

type teamRepository struct {
	db   *sqlx.DB
	pCtx context.Context
}

func TeamRepository(db *sqlx.DB, pCtx context.Context) ITeamRepository {
	return &teamRepository{
		db:   db,
		pCtx: pCtx,
	}
}

// when create team
// insert team
// insert team member (Owner) with role = OWNER (Check if Owner is exist in User table first)
// if have members then loop insert team member with role = MEMBER (Check if Member is exist in User table first)
func (r *teamRepository) CreateTeam(userId string, req *team.CreateTeamReq) (*team.CreateTeamRes, error) {
	res := new(team.CreateTeamRes)
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(r.pCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction create team failed: %v", err)
	}

	queryTeam := `
	INSERT INTO "Team" (
		team_name,
		team_desc,
		team_code
		)
	VALUES ($1, $2, $3)
	RETURNING "team_id", "team_name", "team_poster";
	`
	if err := tx.QueryRowContext(ctx,
		queryTeam,
		req.TeamName,
		req.TeamDesc,
		req.TeamCode,
	).Scan(&res.TeamId, &res.TeamName, &res.TeamPoster); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"Team_team_code_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("team code has been used")
		default:
			return nil, fmt.Errorf("insert team failed: %v", err)
		}
	}

	//insert Permission for team (default allow_task = true, allow_file = true, allow_invite = true)
	queryPermission := `INSERT INTO "Permission" (team_id) VALUES ($1);`
	if _, err := tx.ExecContext(ctx, queryPermission, res.TeamId); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("insert permission failed: %v", err)
	}

	queryTeamOwner := `
	INSERT INTO "TeamMember" (
		team_id,
		user_id,
		role
		)
	VALUES ($1, $2, $3);
	`
	if _, err := tx.ExecContext(ctx,
		queryTeamOwner,
		res.TeamId,
		userId,
		"OWNER",
	); err != nil {
		tx.Rollback()
		switch err.Error() {
		case "ERROR: insert or update on table \"TeamMember\" violates foreign key constraint \"TeamMember_user_id_fkey\" (SQLSTATE 23503)":
			return nil, fmt.Errorf("owner does not exist")
		default:
			return nil, fmt.Errorf("insert team owner failed: %v", err)

		}
	}

	if len(req.Members) > 0 {
		queryTeamMember := `
		INSERT INTO "TeamMember" (
			team_id,
			user_id,
			role
			)
		VALUES ($1, $2, $3);
		`
		for _, member := range req.Members {
			if _, err := tx.ExecContext(ctx,
				queryTeamMember,
				res.TeamId,
				member,
				"MEMBER",
			); err != nil {
				tx.Rollback()
				switch err.Error() {
				case "ERROR: insert or update on table \"TeamMember\" violates foreign key constraint \"TeamMember_user_id_fkey\" (SQLSTATE 23503)":
					return nil, fmt.Errorf("some member does not exist")
				default:
					return nil, fmt.Errorf("insert team member failed: %v", err)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction create team failed: %v", err)
	}

	return res, nil
}

func (r *teamRepository) GetTeamById(teamId string) (*team.GetTeamByIdRes, error) {

	query := `
	SELECT
		"t"."team_id",
		"t"."team_name",
		"t"."team_poster",
		"p"."allow_task",
		"p"."allow_file",
		"p"."allow_invite"
	FROM "Team" "t"
	INNER JOIN "Permission" "p"
	ON "t"."team_id" = "p"."team_id"
	WHERE "t"."team_id" = $1;
	`
	// fmt.Println("teamId", teamId)
	team := new(team.GetTeamByIdRes)
	if err := r.db.Get(team, query, teamId); err != nil {
		return nil, fmt.Errorf("get team by id failed: %v", err)
	}

	return team, nil
}

// when join team
// Check if code team exist then insert team member with role = MEMBER (Check if Member is exist in User table first)
// Check if user already join team then return error
func (r *teamRepository) JoinTeam(req *team.JoinTeamReq) (*team.JoinTeamRes, error) {
	res := new(team.JoinTeamRes)
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	queryCheckCode := `
	SELECT
		team_id,
		team_name,
		team_poster
	FROM "Team"
	WHERE "team_code" = $1;
	`

	if err := r.db.Get(res, queryCheckCode, req.TeamCode); err != nil {
		return nil, fmt.Errorf("team code does not exist")
	}

	queryCheckMmeber := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "TeamMember"
	WHERE "user_id" = $1
	AND "team_id" = $2;
	`

	var isMember bool
	if err := r.db.Get(&isMember, queryCheckMmeber, req.UserId, res.TeamId); err != nil {
		return nil, fmt.Errorf("check member in team failed: %v", err)
	}

	if isMember {
		return nil, fmt.Errorf("user already join team")
	}

	//check if code team is exist then insert team member
	queryTeamMember := `
	INSERT INTO "TeamMember" (
		team_id,
		user_id,
		role
		)
	VALUES ($1, $2, $3);
	`
	if _, err := r.db.ExecContext(ctx,
		queryTeamMember,
		res.TeamId,
		req.UserId,
		"MEMBER",
	); err != nil {
		switch err.Error() {
		case "ERROR: insert or update on table \"TeamMember\" violates foreign key constraint \"TeamMember_user_id_fkey\" (SQLSTATE 23503)":
			return nil, fmt.Errorf("user does not exist")
		default:
			return nil, fmt.Errorf("insert team member failed: %v", err)
		}
	}

	return res, nil
}

func (r *teamRepository) GetTeamByUserId(userId string) ([]*team.GetTeamByUserIdRes, error) {
	query := `
	SELECT
		"t"."team_id",
		"t"."team_name",
		"t"."team_poster"
	FROM "Team" "t"
	INNER JOIN "TeamMember" "tm"
	ON "t"."team_id" = "tm"."team_id"
	WHERE "tm"."user_id" = $1;
	`
	teams := make([]*team.GetTeamByUserIdRes, 0)
	if err := r.db.Select(&teams, query, userId); err != nil {
		return nil, fmt.Errorf("get team by user id failed: %v", err)
	}

	return teams, nil
}

func (r *teamRepository) InviteMember(team_id string, req *team.InviteMemberReq) error {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	// loop insert team member with role = MEMBER (Check if Member is exist in User table first)
	queryTeamMember := `
	INSERT INTO "TeamMember" (
		team_id,
		user_id,
		role
		)
	VALUES ($1, $2, $3);
	`
	for _, member := range req.Users {

		queryCheckMmeber := `
		SELECT
			(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
		FROM "TeamMember"
		WHERE "user_id" = $1
		AND "team_id" = $2;
		`

		var isMember bool
		if err := r.db.Get(&isMember, queryCheckMmeber, member, team_id); err != nil {
			return fmt.Errorf("check member in team failed: %v", err)
		}

		if isMember {
			return fmt.Errorf("user already join team")
		}

		if _, err := r.db.ExecContext(ctx,
			queryTeamMember,
			team_id,
			member,
			"MEMBER",
		); err != nil {
			switch err.Error() {
			case "ERROR: insert or update on table \"TeamMember\" violates foreign key constraint \"TeamMember_user_id_fkey\" (SQLSTATE 23503)":
				return fmt.Errorf("some member does not exist")
			case "ERROR: insert or update on table \"TeamMember\" violates foreign key constraint \"TeamMember_team_id_fkey\" (SQLSTATE 23503)":
				return fmt.Errorf("team does not exist")
			default:
				return fmt.Errorf("invite member failed: %v", err)
			}
		}
	}

	return nil
}

func (r *teamRepository) GetMemberTeam(teamId string) ([]*team.GetMemberTeamRes, error) {
	query := `
	SELECT
		"tm"."member_id",
		"u"."username",
		"u"."email",
		"tm"."role"
	FROM "User" "u"
	INNER JOIN "TeamMember" "tm"
	ON "u"."user_id" = "tm"."user_id"
	WHERE "tm"."team_id" = $1;
	`
	members := make([]*team.GetMemberTeamRes, 0)
	if err := r.db.Select(&members, query, teamId); err != nil {
		return nil, fmt.Errorf("get member team failed: %v", err)
	}

	return members, nil
}

func (r *teamRepository) DeleteMember(memberId string) error {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	//check if member is owner then return error
	queryCheckOwner := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "TeamMember"
	WHERE "member_id" = $1
	AND "role" = 'OWNER';
	`

	var isOwner bool
	if err := r.db.Get(&isOwner, queryCheckOwner, memberId); err != nil {
		return fmt.Errorf("check owner failed: %v", err)
	}

	if isOwner {
		return fmt.Errorf("cannot delete owner")
	}

	query := `
	DELETE FROM "TeamMember"
	WHERE "member_id" = $1;
	`
	if _, err := r.db.ExecContext(ctx, query, memberId); err != nil {
		return fmt.Errorf("delete member failed: %v", err)
	}

	return nil
}

func (r *teamRepository) GetAboutTeam(teamId string) (*team.GetAboutTeamRes, error) {
	query := `
	SELECT
		"team_name",
		"team_desc",
		"team_poster"
	FROM "Team"
	WHERE "team_id" = $1;
	`
	about := new(team.GetAboutTeamRes)
	if err := r.db.Get(about, query, teamId); err != nil {
		return nil, fmt.Errorf("get about team failed: %v", err)
	}

	return about, nil
}

func (r *teamRepository) GetSettingTeam(teamId string) (*team.GetSettingTeamRes, error) {
	query := `
	SELECT
		"t"."team_name",
		"t"."team_desc",
		"t"."team_poster",
		"t"."team_code",
		"p"."allow_task",
		"p"."allow_file",
		"p"."allow_invite"
	FROM "Team" "t"
	INNER JOIN "Permission" "p"
	ON "t"."team_id" = "p"."team_id"
	WHERE "t"."team_id" = $1;
	`
	setting := new(team.GetSettingTeamRes)
	if err := r.db.Get(setting, query, teamId); err != nil {
		return nil, fmt.Errorf("get setting team failed: %v", err)
	}

	return setting, nil
}

func (r *teamRepository) UpdateTeam(teamId string, req *team.UpdateTeamReq) error {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	query := `
	UPDATE "Team" SET`

	if req.TeamName != "" {
		values = append(values, req.TeamName)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"team_name" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.TeamDesc != "" {
		values = append(values, req.TeamDesc)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"team_desc" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.TeamPoster != "" {
		values = append(values, req.TeamPoster)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"team_poster" = $%d?`, lastIndex))

		lastIndex++
	}

	values = append(values, teamId)

	queryClose := fmt.Sprintf(`
	WHERE "team_id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}
	query += queryClose

	if _, err := r.db.ExecContext(ctx, query, values...); err != nil {
		return fmt.Errorf("update profile team failed: %v", err)
	}

	return nil
}

func (r *teamRepository) UpdatePermission(teamId string, req *team.UpdatePermissionReq) error {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
    UPDATE "Permission" SET
        %s = $1
    WHERE "team_id" = $2;
    `, req.PermissionType)

	if _, err := r.db.ExecContext(ctx, query, req.Value, teamId); err != nil {
		return fmt.Errorf("update permission failed: %v", err)
	}

	return nil
}

func (r *teamRepository) UpdateCodeTeam(teamId string, req *team.UpdateCodeTeamReq) error {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	query := `
	UPDATE "Team" SET
		"team_code" = $1
	WHERE "team_id" = $2;
	`
	if _, err := r.db.ExecContext(ctx, query, req.TeamCode, teamId); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"Team_team_code_key\" (SQLSTATE 23505)":
			return fmt.Errorf("team code has been used")
		default:
			return fmt.Errorf("update code team failed: %v", err)
		}
	}

	return nil
}

func (r *teamRepository) DeleteTeam(teamId string) error {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	query := `
	DELETE FROM "Team"
	WHERE "team_id" = $1;
	`
	if _, err := r.db.ExecContext(ctx, query, teamId); err != nil {
		return fmt.Errorf("delete team failed: %v", err)
	}

	return nil
}

func (r *teamRepository) ExitTeam(userId, teamId string) error {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	fmt.Println("userId", userId)
	fmt.Println("teamId", teamId)

	query := `
	DELETE FROM "TeamMember"
	WHERE "user_id" = $1
	AND "team_id" = $2;
	`
	if _, err := r.db.ExecContext(ctx, query, userId, teamId); err != nil {
		return fmt.Errorf("exit team failed: %v", err)
	}

	return nil
}
