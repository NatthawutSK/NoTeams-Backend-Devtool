package middlewaresRepositories

import (
	"github.com/jmoiron/sqlx"
)

type authTeamRes struct {
	IsMember bool `db:"is_member"`
	IsOwner  bool `db:"is_owner"`
}

type IMiddlewaresRepository interface {
	FindAccessToken(userId, accessToken string) bool
	// IsMemberInTeam(userId, teamId string) bool
	// IsOwnerInTeam(userId, teamId string) bool
	IsAllowInviteMember(teamId string) bool
	IsAllowTask(teamId string) bool
	IsAllowFile(teamId string) bool
	AuthTeam(userId, teamId string) (bool, bool)
}

type middlewaresRepository struct {
	db *sqlx.DB
}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresRepository{
		db: db,
	}
}

func (r *middlewaresRepository) FindAccessToken(userId, accessToken string) bool {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "Oauth"
	WHERE "user_id" = $1
	AND "access_token" = $2;`

	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}
	return check
}

// func (r *middlewaresRepository) IsMemberInTeam(userId, teamId string) bool {
// 	query := `
// 	SELECT
// 		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
// 	FROM "TeamMember"
// 	WHERE "user_id" = $1
// 	AND "team_id" = $2;`

// 	var check bool
// 	if err := r.db.Get(&check, query, userId, teamId); err != nil {
// 		return false
// 	}
// 	return check
// }

// func (r *middlewaresRepository) IsOwnerInTeam(userId, teamId string) bool {
// 	query := `
// 	SELECT
// 		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
// 	FROM "TeamMember"
// 	WHERE "user_id" = $1
// 	AND "team_id" = $2
// 	AND "role" = 'OWNER';`

//		var check bool
//		if err := r.db.Get(&check, query, userId, teamId); err != nil {
//			return false
//		}
//		return check
//	}
func (r *middlewaresRepository) AuthTeam(userId, teamId string) (bool, bool) {
	query := `
	SELECT
    	COUNT(*) = 1 AS is_member,
    	COALESCE(BOOL_OR("role" = 'OWNER'), FALSE) AS is_owner
	FROM "TeamMember"
	WHERE "user_id" = $1
	AND "team_id" = $2;`

	authTeam := new(authTeamRes)
	if err := r.db.Get(authTeam, query, userId, teamId); err != nil {
		return false, false
	}

	return authTeam.IsMember, authTeam.IsOwner
}

func (r *middlewaresRepository) IsAllowInviteMember(teamId string) bool {
	query := `
	SELECT
		allow_invite
	FROM "Permission"
	WHERE "team_id" = $1;`

	var check bool
	if err := r.db.Get(&check, query, teamId); err != nil {
		return false
	}
	return check
}

func (r *middlewaresRepository) IsAllowTask(teamId string) bool {
	query := `
	SELECT
		allow_task
	FROM "Permission"
	WHERE "team_id" = $1;`

	var check bool
	if err := r.db.Get(&check, query, teamId); err != nil {
		return false
	}
	return check
}

func (r *middlewaresRepository) IsAllowFile(teamId string) bool {
	query := `
	SELECT
		allow_file
	FROM "Permission"
	WHERE "team_id" = $1;`

	var check bool
	if err := r.db.Get(&check, query, teamId); err != nil {
		return false
	}
	return check
}
