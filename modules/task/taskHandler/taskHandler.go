package taskHandler

import (
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/entities"
	"github.com/NatthawutSK/NoTeams-Backend/modules/task"
	"github.com/NatthawutSK/NoTeams-Backend/modules/task/taskUsecase"
	"github.com/gofiber/fiber/v2"
)

type taskHandlerErrorCode string

const (
	createTaskErr taskHandlerErrorCode = "task-001"
)

type ITaskHandler interface {
	AddTask(c *fiber.Ctx) error
}

type taskHandler struct {
	taskUsecase taskUsecase.ITaskUsecase
}

func TaskHandler(taskUsecase taskUsecase.ITaskUsecase) ITaskHandler {
	return &taskHandler{
		taskUsecase: taskUsecase,
	}
}

func (h *taskHandler) AddTask(c *fiber.Ctx) error {
	req := new(task.AddTaskReq)
	teamId := strings.TrimSpace(c.Params("team_id"))

	//validate request
	validate := entities.ContextWrapper(c)
	if err := validate.BindRi(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(createTaskErr),
			err.Error(),
		).Res()
	}

	res, err := h.taskUsecase.AddTask(teamId, req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(createTaskErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		res,
	).Res()
}
