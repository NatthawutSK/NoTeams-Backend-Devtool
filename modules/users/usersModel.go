package users

type UserRegisterReq struct {
	Email    string `db:"email" json:"email" form:"email" validate:"required,email"`
	Password string `db:"password" json:"password" form:"password" validate:"required,min=8,max=32"`
	Username string `db:"username" json:"username" form:"username" validate:"required,min=4,max=32"`
	Dob      string `db:"dob" json:"dob" form:"dob" validate:"required,datetime=2006-01-02"`
	Phone    string `db:"phone" json:"phone" form:"phone" validate:"required,min=10,max=10,number"`
	Bio      string `db:"bio" json:"bio" form:"bio" validate:"min=0,max=255"`
}

type UserLoginReq struct {
	Email    string `db:"email" json:"email" form:"email" validate:"required,email"`
	Password string `db:"password" json:"password" form:"password" validate:"required,min=8,max=32"`
}

type UserRefreshCredentialReq struct {
	RefreshToken string `db:"refresh_token" json:"refresh_token" form:"refresh_token" validate:"required"`
}

type UserRemoveCredentialReq struct {
	OauthId string `db:"id" json:"oauth_id" form:"oauth_id" validate:"required"`
}

type UserUpdateProfileReq struct {
	Username  string `db:"username" json:"username" form:"username" validate:"omitempty,min=4,max=32"`
	Dob       string `db:"dob" json:"dob" form:"dob" validate:"omitempty,datetime=2006-01-02"`
	Phone     string `db:"phone" json:"phone" form:"phone" validate:"omitempty,min=10,max=10,number"`
	Bio       string `db:"bio" json:"bio" form:"bio" validate:"omitempty,min=0,max=255"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url"`
}
