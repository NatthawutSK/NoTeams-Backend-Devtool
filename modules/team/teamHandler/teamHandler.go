package teamHandler

import (
	"fmt"
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/entities"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team/teamUsecase"
	"github.com/gofiber/fiber/v2"
)

type teamHandlerErrorCode string

const (
	createTeamErr        teamHandlerErrorCode = "team-001"
	getTeamByIdErr       teamHandlerErrorCode = "team-002"
	joinTeamErr          teamHandlerErrorCode = "team-003"
	getTeamByUserIdErr   teamHandlerErrorCode = "team-004"
	inviteMemberErr      teamHandlerErrorCode = "team-005"
	getMemberTeamErr     teamHandlerErrorCode = "team-006"
	deleteMemberErr      teamHandlerErrorCode = "team-007"
	getAboutTeamErr      teamHandlerErrorCode = "team-008"
	getSettingTeamErr    teamHandlerErrorCode = "team-009"
	updateProfileTeamErr teamHandlerErrorCode = "team-010"
	updatePermissionErr  teamHandlerErrorCode = "team-011"
	updateCodeTeamErr    teamHandlerErrorCode = "team-012"
	deleteTeamErr        teamHandlerErrorCode = "team-013"
	exitTeamErr          teamHandlerErrorCode = "team-014"
)

type ITeamHandler interface {
	CreateTeam(c *fiber.Ctx) error
	GetTeamById(c *fiber.Ctx) error
	JoinTeam(c *fiber.Ctx) error
	GetTeamByUserId(c *fiber.Ctx) error
	InviteMember(c *fiber.Ctx) error
	GetMemberTeam(c *fiber.Ctx) error
	DeleteMember(c *fiber.Ctx) error
	GetAboutTeam(c *fiber.Ctx) error
	GetSettingTeam(c *fiber.Ctx) error
	UpdateProfileTeam(c *fiber.Ctx) error
	UpdatePermissionTeam(c *fiber.Ctx) error
	UpdateCodeTeam(c *fiber.Ctx) error
	DeleteTeam(c *fiber.Ctx) error
	ExitTeam(c *fiber.Ctx) error
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
	userId := c.Locals("userId").(string)

	fmt.Println("userId", userId)

	//validate request
	validate := entities.ContextWrapper(c)
	if err := validate.BindRi(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(createTeamErr),
			err.Error(),
		).Res()
	}

	result, err := h.teamUsecase.CreateTeam(userId, req)
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
	role := c.Locals("role").(string)

	result, err := h.teamUsecase.GetTeamById(teamId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getTeamByIdErr),
			err.Error(),
		).Res()
	}

	result.UserRole = role

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
			string(getTeamByUserIdErr),
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

func (h *teamHandler) GetMemberTeam(c *fiber.Ctx) error {
	teamId := strings.TrimSpace(c.Params("team_id"))

	result, err := h.teamUsecase.GetMemberTeam(teamId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getMemberTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		result,
	).Res()
}

func (h *teamHandler) DeleteMember(c *fiber.Ctx) error {
	memberId := strings.TrimSpace(c.Params("member_id"))

	role := c.Locals("role").(string)
	if role != "OWNER" {
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(deleteMemberErr),
			"no permission to delete member",
		).Res()
	}

	err := h.teamUsecase.DeleteMember(memberId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteMemberErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		"delete member success",
	).Res()
}

func (h *teamHandler) GetAboutTeam(c *fiber.Ctx) error {
	teamId := strings.TrimSpace(c.Params("team_id"))

	result, err := h.teamUsecase.GetAboutTeam(teamId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getAboutTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		result,
	).Res()
}

func (h *teamHandler) GetSettingTeam(c *fiber.Ctx) error {
	teamId := strings.TrimSpace(c.Params("team_id"))

	role := c.Locals("role").(string)
	if role != "OWNER" {
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(getSettingTeamErr),
			"no permission to get setting team",
		).Res()
	}

	result, err := h.teamUsecase.GetSettingTeam(teamId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getSettingTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		result,
	).Res()
}

func (h *teamHandler) UpdateProfileTeam(c *fiber.Ctx) error {
	req := new(team.UpdateTeamReq)
	teamId := strings.TrimSpace(c.Params("team_id"))

	role := c.Locals("role").(string)
	if role != "OWNER" {
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(updateProfileTeamErr),
			"no permission to pdate profile team",
		).Res()
	}

	//validate request
	validate := entities.ContextWrapper(c)
	if err := validate.BindRi(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProfileTeamErr),
			err.Error(),
		).Res()
	}

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProfileTeamErr),
			err.Error(),
		).Res()
	}

	avatarFile := form.File["team_poster"]

	if err := h.teamUsecase.UpdateProfileTeam(teamId, req, avatarFile); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProfileTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, "update team profile success").Res()
}

func (h *teamHandler) UpdatePermissionTeam(c *fiber.Ctx) error {
	req := new(team.UpdatePermissionReq)
	teamId := strings.TrimSpace(c.Params("team_id"))

	role := c.Locals("role").(string)
	if role != "OWNER" {
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(updatePermissionErr),
			"no permission to update permission team",
		).Res()
	}

	//validate request
	validate := entities.ContextWrapper(c)
	if err := validate.BindRi(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updatePermissionErr),
			err.Error(),
		).Res()
	}

	if err := h.teamUsecase.UpdatePermission(teamId, req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updatePermissionErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, "update permission success").Res()
}

func (h *teamHandler) UpdateCodeTeam(c *fiber.Ctx) error {
	req := new(team.UpdateCodeTeamReq)
	teamId := strings.TrimSpace(c.Params("team_id"))

	role := c.Locals("role").(string)
	if role != "OWNER" {
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(updateCodeTeamErr),
			"no permission to update code team",
		).Res()
	}

	//validate request
	validate := entities.ContextWrapper(c)
	if err := validate.BindRi(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateCodeTeamErr),
			err.Error(),
		).Res()
	}

	if err := h.teamUsecase.UpdateCodeTeam(teamId, req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateCodeTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, "update code team success").Res()
}

func (h *teamHandler) DeleteTeam(c *fiber.Ctx) error {
	teamId := strings.TrimSpace(c.Params("team_id"))

	role := c.Locals("role").(string)
	if role != "OWNER" {
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(deleteTeamErr),
			"no permission to delete team",
		).Res()
	}

	if err := h.teamUsecase.DeleteTeam(teamId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, "delete team success").Res()
}

func (h *teamHandler) ExitTeam(c *fiber.Ctx) error {
	teamId := strings.TrimSpace(c.Params("team_id"))
	userId := c.Locals("userId").(string)

	role := c.Locals("role").(string)
	if role != "MEMBER" {
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(deleteTeamErr),
			"owner can not exit team",
		).Res()
	}

	if err := h.teamUsecase.ExitTeam(userId, teamId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(exitTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, "exit team success").Res()
}
