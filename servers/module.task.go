package servers

import (
	"context"

	"github.com/NatthawutSK/NoTeams-Backend/modules/task/taskHandler"
	"github.com/NatthawutSK/NoTeams-Backend/modules/task/taskRepository"
	"github.com/NatthawutSK/NoTeams-Backend/modules/task/taskUsecase"
)

type ITaskModule interface {
	Init()
	Repository() taskRepository.ITaskRepository
	Usecase() taskUsecase.ITaskUsecase
	Handler() taskHandler.ITaskHandler
}

type taskModule struct {
	*moduleFactory
	repository taskRepository.ITaskRepository
	usecase    taskUsecase.ITaskUsecase
	handler    taskHandler.ITaskHandler
}

func (m *moduleFactory) TaskModule() ITaskModule {
	ctx := context.Background()
	taskRepository := taskRepository.TaskRepository(m.s.db, ctx)
	taskUsecase := taskUsecase.TaskUsecase(m.s.cfg, taskRepository)
	taskHandler := taskHandler.TaskHandler(taskUsecase)
	return &taskModule{
		moduleFactory: m,
		repository:    taskRepository,
		usecase:       taskUsecase,
		handler:       taskHandler,
	}
}

func (m *taskModule) Init() {
	router := m.r.Group("/task")
	router.Post("/:team_id", m.mid.JwtAuth(), m.mid.AuthTeam(), m.mid.IsAllowTask(), m.handler.AddTask)
}

func (p *taskModule) Repository() taskRepository.ITaskRepository { return p.repository }
func (p *taskModule) Usecase() taskUsecase.ITaskUsecase          { return p.usecase }
func (p *taskModule) Handler() taskHandler.ITaskHandler          { return p.handler }
