package usersRepository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NatthawutSK/NoTeams-Backend/modules/users"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	GetProfile(userId string) (*users.User, error)
	FindOneUserByEmailOrUsername(email, username string) (*users.UserCredentialCheck, error)
	InsertUser(req *users.UserRegisterReq) (IUserRepository, error)
	Result() (*users.User, error)
	InsertOauth(req *users.UserPassport) error
	DeleteOauth(oauthId string) error
	UpdateOauth(req *users.UserToken) error
	FindOneOauth(refreshToken string) (*users.Oauth, error)
	UpdateUserProfile(userId string, req *users.UserUpdateProfileReq) error
}

type usersRepository struct {
	db   *sqlx.DB
	pCtx context.Context
	id   string
}

func UserRepository(db *sqlx.DB, pCtx context.Context) IUserRepository {
	return &usersRepository{
		db:   db,
		pCtx: pCtx,
	}
}

// insert user
func (r *usersRepository) InsertUser(req *users.UserRegisterReq) (IUserRepository, error) {
	ctx, cancel := context.WithTimeout(r.pCtx, 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO "User" (
		email,
		password,
		username,
		dob,
		phone,
		bio
		)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING "user_id";
	`
	if err := r.db.QueryRowContext(ctx,
		query,
		req.Email,
		req.Password,
		req.Username,
		req.Dob,
		req.Phone,
		req.Bio,
	).Scan(&r.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"User_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"User_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		case "ERROR: duplicate key value violates unique constraint \"User_phone_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("phone number has been used")
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
		"u"."user_id",
		"u"."email",
		"u"."username",
		"u"."dob",
		"u"."phone",
		"u"."bio"
	FROM "User" "u"
	WHERE "u"."user_id" = $1
	`

	user := new(users.User)
	if err := r.db.Get(user, query, r.id); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	return user, nil
}

func (r *usersRepository) FindOneUserByEmailOrUsername(email, username string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT
		"user_id",
		"email",
		"password",
		"username",
		"dob",
		"phone",
		"bio",
		"avatar"
	FROM "User"
	WHERE "email" = $1 OR "username" = $2;`
	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email, username); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *usersRepository) FindOneOauth(refreshToken string) (*users.Oauth, error) {
	query := `
	SELECT
		"oauth_id",
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
	ctx, cancel := context.WithTimeout(r.pCtx, 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "Oauth" (
		"user_id",
		"refresh_token",
		"access_token"
	)
	VALUES ($1, $2, $3)
		RETURNING "oauth_id";`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.UserId,
		req.Token.RefreshToken,
		req.Token.AccessToken,
	).Scan(&req.Token.OauthId); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) UpdateOauth(req *users.UserToken) error {
	query := `
	UPDATE "Oauth" SET
		"access_token" = :access_token,
		"refresh_token" = :refresh_token
	WHERE "oauth_id" = :oauth_id;`

	if _, err := r.db.NamedExecContext(r.pCtx, query, req); err != nil {
		return fmt.Errorf("update oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) DeleteOauth(oauthId string) error {
	query := `
	DELETE FROM "Oauth"
	WHERE "oauth_id" = $1;`

	if _, err := r.db.ExecContext(r.pCtx, query, oauthId); err != nil {
		return fmt.Errorf("oauth not found")
	}
	return nil
}

func (r *usersRepository) GetProfile(userId string) (*users.User, error) {
	query := `
	SELECT
		"user_id",
		"email",
		"username",
		"dob",
		"phone",
		"bio",
		"avatar"
	FROM "User"
	WHERE "user_id" = $1;`

	profile := new(users.User)
	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}
	return profile, nil
}

func (r *usersRepository) UpdateUserProfile(userId string, req *users.UserUpdateProfileReq) error {
	ctx, cancel := context.WithTimeout(r.pCtx, 10*time.Second)
	defer cancel()

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	query := `
	UPDATE "User" SET`

	if req.Username != "" {
		values = append(values, req.Username)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"username" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.Phone != "" {
		values = append(values, req.Phone)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"phone" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.Dob != "" {
		values = append(values, req.Dob)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"dob" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.Bio != "" {
		values = append(values, req.Bio)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"bio" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.AvatarUrl != "" {
		values = append(values, req.AvatarUrl)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"avatar" = $%d?`, lastIndex))

		lastIndex++
	}

	values = append(values, userId)

	queryClose := fmt.Sprintf(`
	WHERE "user_id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}
	query += queryClose

	if _, err := r.db.ExecContext(ctx, query, values...); err != nil {
		return fmt.Errorf("update profile user failed: %v", err)
	}

	return nil
}
