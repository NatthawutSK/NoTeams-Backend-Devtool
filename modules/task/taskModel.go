package task

type AddTaskReq struct {
	TaskName     string `json:"task_name" form:"task_name" validate:"required,min=1,max=32" db:"task_name"`
	TaskDesc     string `json:"task_desc" form:"task_desc" validate:"min=0,max=255" db:"task_desc"`
	TaskStatus   string `json:"task_status" form:"task_status" validate:"required" db:"task_status"`
	TaskDeadline string `json:"task_deadline" form:"task_deadline" db:"task_deadline" validate:"omitempty,datetime=2006-01-02"`
	UserId       string `json:"user_id" form:"user_id" db:"user_id"`
}

type UpdateTaskReq struct {
	TaskId       string `json:"task_id" form:"task_id" validate:"required" db:"task_id"`
	TaskName     string `json:"task_name" form:"task_name" validate:"omitempty,min=1,max=32" db:"task_name"`
	TaskDesc     string `json:"task_desc" form:"task_desc" validate:"omitempty,min=0,max=255" db:"task_desc"`
	TaskDeadline string `json:"task_deadline" form:"task_deadline" db:"task_deadline"`
	UserId       string `json:"user_id" form:"user_id" db:"user_id"`
}

type DeleteTaskReq struct {
	TaskId string `json:"task_id" form:"task_id" validate:"required" db:"task_id"`
}

type MoveTaskReq struct {
	TaskId     string `json:"task_id" form:"task_id" validate:"required" db:"task_id"`
	TaskStatus string `json:"task_status" form:"task_status" validate:"required" db:"task_status"`
}
