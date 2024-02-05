package servers

import (
	"github.com/NatthawutSK/NoTeams-Backend/modules/middlewares/middlewaresHandlers"
	"github.com/NatthawutSK/NoTeams-Backend/modules/middlewares/middlewaresRepositories"
	"github.com/NatthawutSK/NoTeams-Backend/modules/middlewares/middlewaresUsecases"
	"github.com/NatthawutSK/NoTeams-Backend/modules/monitor/monitorHandlers"
	"github.com/NatthawutSK/NoTeams-Backend/modules/users/usersHandlers"
	"github.com/NatthawutSK/NoTeams-Backend/modules/users/usersRepositories"
	"github.com/NatthawutSK/NoTeams-Backend/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func (m *moduleFactory) MonitorModule() {
	handle := monitorHandlers.MonitorHandler(m.s.cfg)
	m.r.Get("/", handle.HealthCheck)
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UserRepository(m.s.db)
	usecase := usersUsecases.UserUsecase(repository, m.s.cfg)
	handler := usersHandlers.UsersHandler(usecase, m.s.cfg)

	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUp)
	router.Post("/signin", handler.SignIn)
	router.Post("/signout", handler.SignOut)
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
}
