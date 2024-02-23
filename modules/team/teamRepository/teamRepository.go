package teamRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/NatthawutSK/NoTeams-Backend/modules/team"
	"github.com/jmoiron/sqlx"
)

type ITeamRepository interface {
	CreateTeam(req *team.CreateTeamReq) (*team.CreateTeamRes, error)
	GetTeamById(teamId string) (*team.GetTeamByIdRes, error)
	JoinTeam(req *team.JoinTeamReq) (*team.JoinTeamRes, error)
	GetTeamByUserId(userId string) ([]*team.GetTeamByUserIdRes, error)
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
func (r *teamRepository) CreateTeam(req *team.CreateTeamReq) (*team.CreateTeamRes, error) {
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
		req.OwnerId,
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
		"team_id",
		"team_name",
		"team_poster"
	FROM "Team"
	WHERE "team_id" = $1;
	`
	fmt.Println("teamId", teamId)
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
