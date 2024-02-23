package teamHandler

import (
	"github.com/NatthawutSK/NoTeams-Backend/entities"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team/teamUsecase"
	"github.com/gofiber/fiber/v2"
)

type teamHandlerErrorCode string

const (
	createTeamErr teamHandlerErrorCode = "team-001"
)

type ITeamHandler interface {
	CreateTeam(c *fiber.Ctx) error
}

type teamHandler struct {
	teamUsecase teamUsecase.ITeamUsecase
}

func TeamHandler(teamUsecase teamUsecase.ITeamUsecase) ITeamHandler {
	return &teamHandler{
		teamUsecase: teamUsecase,
	}
}

func (h *teamHandler) CreateTeam(c *fiber.Ctx) error {
	req := new(team.CreateTeamReq)

	//validate request
	validate := entities.ContextWrapper(c)
	if err := validate.BindRi(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(createTeamErr),
			err.Error(),
		).Res()
	}

	result, err := h.teamUsecase.CreateTeam(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(createTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		result,
	).Res()

}
