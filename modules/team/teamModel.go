package team

type CreateTeamReq struct {
	OwnerId  string    `json:"owner_id" form:"owner_id" validate:"required" db:"owner_id"`
	TeamName string    `json:"team_name" form:"team_name" validate:"required,min=1,max=32" db:"team_name"`
	TeamDesc string    `json:"team_desc" form:"team_desc" validate:"min=0,max=255" db:"team_desc"`
	TeamCode string    `json:"team_code" form:"team_code" validate:"required,min=4,max=32" db:"team_code"`
	Members  []*string `json:"members" form:"members"`
}

type JoinTeamReq struct {
	TeamCode string `json:"team_code" form:"team_code" validate:"required,min=4,max=32" db:"team_code"`
	UserId   string `json:"user_id" form:"user_id" validate:"required" db:"user_id"`
}
