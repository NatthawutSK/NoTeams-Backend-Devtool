package users

type UserRegisterReq struct {
	Email    string `db:"email" json:"email" form:"email" validate:"required,email"`
	Password string `db:"password" json:"password" form:"password" validate:"required,min=8,max=32"`
	Username string `db:"username" json:"username" form:"username" validate:"required,min=4,max=32"`
	Dob      string `db:"dob" json:"dob" form:"dob" validate:"required"`
	Phone    string `db:"phone" json:"phone" form:"phone" validate:"required,min=10,max=10"`
	Bio      string `db:"bio" json:"bio" form:"bio" validate:"min=0,max=255"`
}

type UserLoginReq struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserRefreshCredentialReq struct {
	RefreshToken string `db:"refresh_token" json:"refresh_token" form:"refresh_token"`
}

type UserRemoveCredentialReq struct {
	OauthId string `db:"id" json:"oauth_id" form:"oauth_id"`
}
