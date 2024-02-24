package servers

import (
	"context"

	"github.com/NatthawutSK/NoTeams-Backend/modules/team/teamHandler"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team/teamRepository"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team/teamUsecase"
)

type ITeamModule interface {
	Init()
	Repository() teamRepository.ITeamRepository
	Usecase() teamUsecase.ITeamUsecase
	Handler() teamHandler.ITeamHandler
}

type teamModule struct {
	*moduleFactory
	repository teamRepository.ITeamRepository
	usecase    teamUsecase.ITeamUsecase
	handler    teamHandler.ITeamHandler
}

func (m *moduleFactory) TeamModule() ITeamModule {
	ctx := context.Background()
	teamRepository := teamRepository.TeamRepository(m.s.db, ctx)
	teamUsecase := teamUsecase.TeamUsecase(teamRepository, m.s.cfg)
	teamHandler := teamHandler.TeamHandler(teamUsecase)
	return &teamModule{
		moduleFactory: m,
		repository:    teamRepository,
		usecase:       teamUsecase,
		handler:       teamHandler,
	}
}

func (m *teamModule) Init() {
	router := m.r.Group("/teams")

	router.Post("/", m.mid.JwtAuth(), m.handler.CreateTeam)
	router.Get("/:team_id", m.mid.JwtAuth(), m.mid.AuthTeam(), m.handler.GetTeamById)
	router.Post("/join", m.mid.JwtAuth(), m.handler.JoinTeam)
	router.Get("/user/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.GetTeamByUserId)
	router.Post("/invite/:team_id", m.mid.JwtAuth(), m.mid.AuthTeam(), m.mid.IsAllowInvite(), m.handler.InviteMember)
	router.Get("/member/:team_id", m.mid.JwtAuth(), m.mid.AuthTeam(), m.handler.GetMemberTeam)
	router.Delete("/:team_id/member/:member_id", m.mid.JwtAuth(), m.mid.AuthTeam(), m.handler.DeleteMember)
	router.Get("/about/:team_id", m.mid.JwtAuth(), m.mid.AuthTeam(), m.handler.GetAboutTeam)
	router.Get("/setting/:team_id", m.mid.JwtAuth(), m.mid.AuthTeam(), m.handler.GetSettingTeam)
	router.Put("/profile/:team_id", m.mid.JwtAuth(), m.mid.AuthTeam(), m.handler.UpdateProfileTeam)

}

func (p *teamModule) Repository() teamRepository.ITeamRepository { return p.repository }
func (p *teamModule) Usecase() teamUsecase.ITeamUsecase          { return p.usecase }
func (p *teamModule) Handler() teamHandler.ITeamHandler          { return p.handler }
