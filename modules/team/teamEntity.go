package team

type CreateTeamRes struct {
	TeamId     string `json:"team_id" db:"team_id"`
	TeamName   string `json:"team_name" db:"team_name"`
	TeamPoster string `json:"team_poster" db:"team_poster"`
}

type GetTeamByIdRes struct {
	TeamId      string `json:"team_id" db:"team_id"`
	TeamName    string `json:"team_name" db:"team_name"`
	TeamPoster  string `json:"team_poster" db:"team_poster"`
	UserRole    string `json:"user_role" db:"user_role"`
	AllowTask   bool   `json:"allow_task" db:"allow_task"`
	AllowFile   bool   `json:"allow_file" db:"allow_file"`
	AllowInvite bool   `json:"allow_invite" db:"allow_invite"`
}

type JoinTeamRes struct {
	TeamId     string `json:"team_id" db:"team_id"`
	TeamName   string `json:"team_name" db:"team_name"`
	TeamPoster string `json:"team_poster" db:"team_poster"`
}

type GetTeamByUserIdRes struct {
	TeamId     string `json:"team_id" db:"team_id"`
	TeamName   string `json:"team_name" db:"team_name"`
	TeamPoster string `json:"team_poster" db:"team_poster"`
}

type GetMemberTeamRes struct {
	MemberId string `json:"member_id" db:"member_id"`
	UserId   string `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Role     string `json:"role" db:"role"`
}

type GetAboutTeamRes struct {
	TeamName   string `json:"team_name" db:"team_name"`
	TeamDesc   string `json:"team_desc" db:"team_desc"`
	TeamPoster string `json:"team_poster" db:"team_poster"`
}

type GetSettingTeamRes struct {
	TeamName    string `json:"team_name" db:"team_name"`
	TeamDesc    string `json:"team_desc" db:"team_desc"`
	TeamPoster  string `json:"team_poster" db:"team_poster"`
	TeamCode    string `json:"team_code" db:"team_code"`
	AllowTask   bool   `json:"allow_task" db:"allow_task"`
	AllowFile   bool   `json:"allow_file" db:"allow_file"`
	AllowInvite bool   `json:"allow_invite" db:"allow_invite"`
}
