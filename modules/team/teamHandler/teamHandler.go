package teamHandler

import (
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/entities"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team/teamUsecase"
	"github.com/gofiber/fiber/v2"
)

type teamHandlerErrorCode string

const (
	createTeamErr   teamHandlerErrorCode = "team-001"
	getTeamById     teamHandlerErrorCode = "team-002"
	joinTeamErr     teamHandlerErrorCode = "team-003"
	getTeamByUserId teamHandlerErrorCode = "team-004"
	inviteMemberErr teamHandlerErrorCode = "team-005"
)

type ITeamHandler interface {
	CreateTeam(c *fiber.Ctx) error
	GetTeamById(c *fiber.Ctx) error
	JoinTeam(c *fiber.Ctx) error
	GetTeamByUserId(c *fiber.Ctx) error
	InviteMember(c *fiber.Ctx) error
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

func (h *teamHandler) GetTeamById(c *fiber.Ctx) error {
	teamId := strings.TrimSpace(c.Params("team_id"))

	result, err := h.teamUsecase.GetTeamById(teamId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getTeamById),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		result,
	).Res()
}

func (h *teamHandler) JoinTeam(c *fiber.Ctx) error {
	req := new(team.JoinTeamReq)

	//validate request
	validate := entities.ContextWrapper(c)
	if err := validate.BindRi(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(joinTeamErr),
			err.Error(),
		).Res()
	}

	result, err := h.teamUsecase.JoinTeam(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(joinTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		result,
	).Res()
}

func (h *teamHandler) GetTeamByUserId(c *fiber.Ctx) error {
	userId := strings.TrimSpace(c.Params("user_id"))

	result, err := h.teamUsecase.GetTeamByUserId(userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getTeamByUserId),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		result,
	).Res()
}

func (h *teamHandler) InviteMember(c *fiber.Ctx) error {
	team_id := strings.TrimSpace(c.Params("team_id"))
	req := new(team.InviteMemberReq)

	//validate request
	validate := entities.ContextWrapper(c)
	if err := validate.BindRi(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(inviteMemberErr),
			err.Error(),
		).Res()
	}

	err := h.teamUsecase.InviteMember(team_id, req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(inviteMemberErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		"invite member success",
	).Res()
}
