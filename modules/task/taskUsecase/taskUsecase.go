package taskUsecase

import (
	"fmt"
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/modules/task"
	"github.com/NatthawutSK/NoTeams-Backend/modules/task/taskRepository"
)

type ITaskUsecase interface {
	AddTask(teamId string, req *task.AddTaskReq) (*task.AddTaskRes, error)
	UpdateTask(teamId string, req *task.UpdateTaskReq) error
}

type taskUsecase struct {
	taskRepo taskRepository.ITaskRepository
	cfg      config.IConfig
}

func TaskUsecase(cfg config.IConfig, taskRepo taskRepository.ITaskRepository) ITaskUsecase {
	return &taskUsecase{
		taskRepo: taskRepo,
		cfg:      cfg,
	}
}

func (u *taskUsecase) AddTask(teamId string, req *task.AddTaskReq) (*task.AddTaskRes, error) {
	req.TaskStatus = strings.ToUpper(req.TaskStatus)

	status := map[string]bool{
		"TODO":  true,
		"DOING": true,
		"DONE":  true,
	}

	if !status[req.TaskStatus] {
		return nil, fmt.Errorf("invalid task status")
	}

	res, err := u.taskRepo.AddTask(teamId, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *taskUsecase) UpdateTask(teamId string, req *task.UpdateTaskReq) error {

	if err := u.taskRepo.UpdateTask(teamId, req); err != nil {
		return err
	}

	return nil
}
