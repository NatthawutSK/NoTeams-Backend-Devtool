package servers

import (
	"context"

	"github.com/NatthawutSK/NoTeams-Backend/modules/users/usersHandler"
	"github.com/NatthawutSK/NoTeams-Backend/modules/users/usersRepository"
	"github.com/NatthawutSK/NoTeams-Backend/modules/users/usersUsecase"
)

type IUserModule interface {
	Init()
	Repository() usersRepository.IUserRepository
	Usecase() usersUsecase.IUserUsecase
	Handler() usersHandler.IUsersHandler
}

type userModule struct {
	*moduleFactory
	repository usersRepository.IUserRepository
	usecase    usersUsecase.IUserUsecase
	handler    usersHandler.IUsersHandler
}

func (m *moduleFactory) UserModule() IUserModule {
	ctx := context.Background()
	// fileUsecase := filesUsecase.FilesUsecase(m.s.cfg)
	userRepository := usersRepository.UserRepository(m.s.db, ctx)
	userUsecase := usersUsecase.UserUsecase(userRepository, m.s.cfg)
	userHandler := usersHandler.UsersHandler(userUsecase, m.s.cfg)
	return &userModule{
		moduleFactory: m,
		repository:    userRepository,
		usecase:       userUsecase,
		handler:       userHandler,
	}
}

func (m *userModule) Init() {
	router := m.r.Group("/users")

	router.Post("/signup", m.handler.SignUp)
	router.Post("/signin", m.handler.SignIn)
	router.Post("/refresh", m.mid.JwtAuth(), m.handler.RefreshPassport)
	router.Post("/signout", m.mid.JwtAuth(), m.handler.SignOut)
	router.Get("/profile/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.GetUserProfile)
	router.Get("/find", m.mid.JwtAuth(), m.handler.FindOneUserByEmailOrUsername)
	router.Put("/profile/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.UpdateUserProfile)

}

func (p *userModule) Repository() usersRepository.IUserRepository { return p.repository }
func (p *userModule) Usecase() usersUsecase.IUserUsecase          { return p.usecase }
func (p *userModule) Handler() usersHandler.IUsersHandler         { return p.handler }
