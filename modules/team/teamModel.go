package team

type CreateTeamReq struct {
	TeamName string    `json:"team_name" form:"team_name" validate:"required,min=1,max=32" db:"team_name"`
	TeamDesc string    `json:"team_desc" form:"team_desc" validate:"min=0,max=255" db:"team_desc"`
	TeamCode string    `json:"team_code" form:"team_code" validate:"required,min=4,max=32" db:"team_code"`
	Members  []*string `json:"members" form:"members"`
}

type JoinTeamReq struct {
	TeamCode string `json:"team_code" form:"team_code" validate:"required,min=4,max=32" db:"team_code"`
	UserId   string `json:"user_id" form:"user_id" validate:"required" db:"user_id"`
}

type InviteMemberReq struct {
	Users []*string `json:"users" form:"users" validate:"required" db:"users"`
}

type UpdateTeamReq struct {
	TeamName   string `json:"team_name" form:"team_name" validate:"omitempty,min=1,max=32" db:"team_name"`
	TeamDesc   string `json:"team_desc" form:"team_desc" validate:"omitempty,min=0,max=255" db:"team_desc"`
	TeamPoster string `json:"team_poster" form:"team_poster" db:"team_poster"`
}

type UpdatePermissionReq struct {
	PermissionType string `json:"permission_type" form:"permission_type" validate:"required"`
	Value          bool   `json:"value" form:"value" validate:"boolean"`
}

type UpdateCodeTeamReq struct {
	TeamCode string `json:"team_code" form:"team_code" validate:"required,min=4,max=32" db:"team_code"`
}
