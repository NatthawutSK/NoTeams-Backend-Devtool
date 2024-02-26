package taskRepository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NatthawutSK/NoTeams-Backend/modules/task"
	"github.com/jmoiron/sqlx"
)

type ITaskRepository interface {
	AddTask(teamId string, req *task.AddTaskReq) (*task.AddTaskRes, error)
	UpdateTask(teamId string, req *task.UpdateTaskReq) error
}

type taskRepository struct {
	db   *sqlx.DB
	pCtx context.Context
}

func TaskRepository(db *sqlx.DB, pCtx context.Context) ITaskRepository {
	return &taskRepository{
		db:   db,
		pCtx: pCtx,
	}
}

func (r *taskRepository) AddTask(teamId string, req *task.AddTaskReq) (*task.AddTaskRes, error) {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	var userId *string

	// Check if req.UserId is not an empty string
	if req.UserId != "" {
		// If req.UserId is not an empty string, assign its address to userId
		userId = &req.UserId

		//check user is in team or not
		query := `
		SELECT
			(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
		FROM "TeamMember"
		WHERE "user_id" = $1
		AND "team_id" = $2;`

		var check bool
		if err := r.db.Get(&check, query, userId, teamId); err != nil {
			return nil, fmt.Errorf("check user in team failed: %v", err)
		}

		if !check {
			return nil, fmt.Errorf("user not in team")
		}

	}

	query := `INSERT INTO "Task" (
				task_name,
				task_desc,
		  		task_status,
		   		task_deadline,
				user_id,
				team_id
		    )
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING "task_id";`

	taskRes := new(task.AddTaskRes)
	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.TaskName,
		req.TaskDesc,
		req.TaskStatus,
		req.TaskDeadline,
		userId,
		teamId,
	).Scan(&taskRes.TaskId); err != nil {
		switch err.Error() {
		case "ERROR: insert or update on table \"Task\" violates foreign key constraint \"Task_user_id_fkey\" (SQLSTATE 23503)":
			return nil, fmt.Errorf("user not found")
		case "ERROR: insert or update on table \"Task\" violates foreign key constraint \"Task_team_id_fkey\" (SQLSTATE 23503)":
			return nil, fmt.Errorf("team not found")
		default:
			return nil, fmt.Errorf("add task failed: %v", err)
		}
	}

	return taskRes, nil
}

func (r *taskRepository) UpdateTask(teamId string, req *task.UpdateTaskReq) error {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	query := `
	UPDATE "Task" SET`

	if req.TaskName != "" {
		values = append(values, req.TaskName)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"task_name" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.TaskDesc != "" {
		values = append(values, req.TaskDesc)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"task_desc" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.TaskDeadline != "" {
		values = append(values, req.TaskDeadline)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"task_deadline" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.UserId != "" {
		//check user is in team or not
		query := `
		SELECT
			(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
		FROM "TeamMember"
		WHERE "user_id" = $1
		AND "team_id" = $2;`

		var check bool
		if err := r.db.Get(&check, query, req.UserId, teamId); err != nil {
			return fmt.Errorf("check user in team failed: %v", err)
		}

		if !check {
			return fmt.Errorf("user not in team")
		}

		values = append(values, req.UserId)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"user_id" = $%d?`, lastIndex))

		lastIndex++
	}

	values = append(values, req.TaskId)

	queryClose := fmt.Sprintf(`
	WHERE "task_id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}
	query += queryClose

	if _, err := r.db.ExecContext(ctx, query, values...); err != nil {
		return fmt.Errorf("update task failed: %v", err)
	}

	return nil
}
