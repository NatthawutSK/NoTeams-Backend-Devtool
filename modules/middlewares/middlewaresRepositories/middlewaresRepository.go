package middlewaresRepositories

import (
	"github.com/jmoiron/sqlx"
)

type IMiddlewaresRepository interface {
	FindAccessToken(userId, accessToken string) bool
	CheckMemberInTeam(userId, teamId string) bool
	CheckOwnerInTeam(userId, teamId string) bool
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

func (r *middlewaresRepository) CheckMemberInTeam(userId, teamId string) bool {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "TeamMember"
	WHERE "user_id" = $1
	AND "team_id" = $2;`

	var check bool
	if err := r.db.Get(&check, query, userId, teamId); err != nil {
		return false
	}
	return check
}

func (r *middlewaresRepository) CheckOwnerInTeam(userId, teamId string) bool {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "TeamMember"
	WHERE "user_id" = $1
	AND "team_id" = $2
	AND "role" = 'OWNER';`

	var check bool
	if err := r.db.Get(&check, query, userId, teamId); err != nil {
		return false
	}
	return check
}
