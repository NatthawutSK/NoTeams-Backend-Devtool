package users

type User struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	Dob      string `db:"dob" json:"dob"`
	Phone    string `db:"phone" json:"phone"`
	Avatar   string `db:"avatar" json:"avatar"`
	Bio      string `db:"bio" json:"bio"`
}

type UserCredentialCheck struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Username string `db:"username" json:"username"`
	Dob      string `db:"dob" json:"dob"`
	Phone    string `db:"phone" json:"phone"`
	Bio      string `db:"bio" json:"bio"`
	Avatar   string `db:"avatar" json:"avatar"`
}

type UserPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           string `db:"id" json:"id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type UserClaims struct {
	Id string `json:"id" db:"id"`
}

type Oauth struct {
	Id     string `db:"id" json:"id"`
	UserId string `db:"user_id" json:"user_id"`
}
