package usersRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/NatthawutSK/NoTeams-Backend/modules/users"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	GetProfile(userId string) (*users.User, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertUser(req *users.UserRegisterReq) (IUserRepository, error)
	Result() (*users.User, error)
	InsertOauth(req *users.UserPassport) error
	DeleteOauth(oauthId string) error
	UpdateOauth(req *users.UserToken) error
	FindOneOauth(refreshToken string) (*users.Oauth, error)
}

type usersRepository struct {
	db *sqlx.DB
	id string
}

func UserRepository(db *sqlx.DB) IUserRepository {
	return &usersRepository{
		db: db,
	}
}

// insert user
func (r *usersRepository) InsertUser(req *users.UserRegisterReq) (IUserRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO "User" (
		email,
		password,
		username
		)
	VALUES ($1, $2, $3)
	RETURNING "id";
	`
	if err := r.db.QueryRowContext(ctx,
		query,
		req.Email,
		req.Password,
		req.Username,
	).Scan(&r.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"User_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"User_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}
	return r, nil
}

// result from insert user
func (r *usersRepository) Result() (*users.User, error) {
	query := `
	SELECT
		"u"."id",
		"u"."email",
		"u"."username"
	FROM "User" "u"
	WHERE "u"."id" = $1
	`

	user := new(users.User)
	if err := r.db.Get(user, query, r.id); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	return user, nil
}

func (r *usersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT
		"id",
		"email",
		"password",
		"username"
	FROM "User"
	WHERE "email" = $1;`
	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *usersRepository) FindOneOauth(refreshToken string) (*users.Oauth, error) {
	query := `
	SELECT
		"id",
		"user_id"
	FROM "Oauth"
	WHERE "refresh_token" = $1;`

	oauth := new(users.Oauth)
	if err := r.db.Get(oauth, query, refreshToken); err != nil {
		return nil, fmt.Errorf("oauth not found")
	}
	return oauth, nil
}

func (r *usersRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "Oauth" (
		"user_id",
		"refresh_token",
		"access_token"
	)
	VALUES ($1, $2, $3)
		RETURNING "id";`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.Id,
		req.Token.RefreshToken,
		req.Token.AccessToken,
	).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) UpdateOauth(req *users.UserToken) error {
	query := `
	UPDATE "Oauth" SET
		"access_token" = :access_token,
		"refresh_token" = :refresh_token
	WHERE "id" = :id;`

	if _, err := r.db.NamedExecContext(context.Background(), query, req); err != nil {
		return fmt.Errorf("update oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) DeleteOauth(oauthId string) error {
	query := `
	DELETE FROM "Oauth"
	WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, oauthId); err != nil {
		return fmt.Errorf("oauth not found")
	}
	return nil
}

func (r *usersRepository) GetProfile(userId string) (*users.User, error) {
	query := `
	SELECT
		"id",
		"email",
		"username"
	FROM "User"
	WHERE "id" = $1;`

	profile := new(users.User)
	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}
	return profile, nil
}
