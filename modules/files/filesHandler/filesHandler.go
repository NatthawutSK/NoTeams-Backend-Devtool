package filesHandler

import (
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/entities"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files/filesUsecase"
	"github.com/gofiber/fiber/v2"
)

type FileHandlerErrCode string

const (
	uploadFilesErr  FileHandlerErrCode = "files-001"
	getFilesTeamErr FileHandlerErrCode = "files-002"
	deleteFilesErr  FileHandlerErrCode = "files-003"
)

type IFileHandler interface {
	GetFilesTeam(c *fiber.Ctx) error
	UploadFilesTeam(c *fiber.Ctx) error
	DeleteFilesTeam(c *fiber.Ctx) error
}

type fileHandler struct {
	cfg         config.IConfig
	fileUsecase filesUsecase.IFilesUsecase
}

func FileHandler(cfg config.IConfig, fileUsecase filesUsecase.IFilesUsecase) IFileHandler {
	return &fileHandler{
		cfg:         cfg,
		fileUsecase: fileUsecase,
	}
}

func (h *fileHandler) UploadFilesTeam(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	teamId := strings.TrimSpace(c.Params("team_id"))

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFilesErr),
			err.Error(),
		).Res()
	}

	filesReq := form.File["files"]
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFilesErr),
			err.Error(),
		).Res()
	}

	if len(filesReq) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFilesErr),
			"no files found",
		).Res()
	}

	res, err := h.fileUsecase.UploadFilesTeam(userId, teamId, filesReq)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFilesErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		res,
	).Res()

}

func (h *fileHandler) GetFilesTeam(c *fiber.Ctx) error {
	teamId := strings.TrimSpace(c.Params("team_id"))
	filesTeam, err := h.fileUsecase.GetFilesTeam(teamId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getFilesTeamErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		filesTeam,
	).Res()
}

func (h *fileHandler) DeleteFilesTeam(c *fiber.Ctx) error {

	req := new(files.DeleteFilesTeamReq)

	//validate request
	validate := entities.ContextWrapper(c)
	if err := validate.BindRi(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteFilesErr),
			err.Error(),
		).Res()
	}

	err := h.fileUsecase.DeleteFilesTeam(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteFilesErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		"deleted file success",
	).Res()
}
