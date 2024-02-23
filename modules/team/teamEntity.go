package team

type CreateTeamRes struct {
	TeamId     string `json:"team_id" db:"team_id"`
	TeamName   string `json:"team_name" db:"team_name"`
	TeamPoster string `json:"team_poster" db:"team_poster"`
}
