package team

type CreateTeamRes struct {
	TeamId     string `json:"team_id" db:"team_id"`
	TeamName   string `json:"team_name" db:"team_name"`
	TeamPoster string `json:"team_poster" db:"team_poster"`
}

type GetTeamByIdRes struct {
	TeamId     string `json:"team_id" db:"team_id"`
	TeamName   string `json:"team_name" db:"team_name"`
	TeamPoster string `json:"team_poster" db:"team_poster"`
	UserRole   string `json:"user_role" db:"user_role"`
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
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Role     string `json:"role" db:"role"`
}

type GetAboutTeamRes struct {
	TeamName   string `json:"team_name" db:"team_name"`
	TeamDesc   string `json:"team_desc" db:"team_desc"`
	TeamPoster string `json:"team_poster" db:"team_poster"`
}
