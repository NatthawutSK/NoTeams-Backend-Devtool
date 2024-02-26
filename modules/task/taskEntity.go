package task

type AddTaskRes struct {
	TaskId string `json:"task_id" db:"task_id"`
}

type GetTaskTeamRes []*GetTaskTeam

type GetTaskTeam struct {
	TaskId       string `json:"task_id" db:"task_id"`
	TaskName     string `json:"task_name" db:"task_name"`
	TaskDesc     string `json:"task_desc" db:"task_desc"`
	TaskStatus   string `json:"task_status" db:"task_status"`
	TaskDeadline string `json:"task_deadline" db:"task_deadline"`
	Username     string `json:"username" db:"username"`
	UserId       string `json:"user_id" db:"user_id"`
}
