package filesHandler

import (
	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/entities"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files/filesUsecase"
	"github.com/gofiber/fiber/v2"
)

type FileHandlerErrCode string

const (
	uploadFilesErr FileHandlerErrCode = "files-001"
)

type IFileHandler interface {
	UploadFiles(c *fiber.Ctx) error
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

// func (h *fileHandler) UploadFiles(c *fiber.Ctx) error {

// 	// form, err := c.MultipartForm()
// 	// if err != nil {
// 	// 	return entities.NewResponse(c).Error(
// 	// 		fiber.ErrBadRequest.Code,
// 	// 		string(uploadFilesErr),
// 	// 		err.Error(),
// 	// 	).Res()
// 	// }
// 	s3Client := s3Conn.S3Connect(h.cfg.S3())

// 	filesReq, err := c.FormFile("files")
// 	if err != nil {
// 		return entities.NewResponse(c).Error(
// 			fiber.ErrBadRequest.Code,
// 			string(uploadFilesErr),
// 			err.Error(),
// 		).Res()
// 	}

// 	// if len(filesReq) == 0 {
// 	// 	return entities.NewResponse(c).Error(
// 	// 		fiber.ErrBadRequest.Code,
// 	// 		string(uploadFilesErr),
// 	// 		"no files found",
// 	// 	).Res()
// 	// }

// 	// res, err := h.fileUsecase.UploadFiles(filesReq)
// 	// Upload the file to S3
// 	url, err := h.fileUsecase.UploadFile(s3Client, h.cfg.S3().S3Bucket(), filesReq.Filename, filesReq)
// 	if err != nil {
// 		return entities.NewResponse(c).Error(
// 			fiber.ErrBadRequest.Code,
// 			string(uploadFilesErr),
// 			err.Error(),
// 		).Res()
// 	}

// 	return entities.NewResponse(c).Success(
// 		fiber.StatusOK,
// 		url,
// 	).Res()

// }

func (h *fileHandler) UploadFiles(c *fiber.Ctx) error {

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFilesErr),
			err.Error(),
		).Res()
	}
	// s3Client := s3Conn.S3Connect(h.cfg.S3())

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

	// Upload the file to S3
	// url, err := h.fileUsecase.UploadFile(s3Client, h.cfg.S3().S3Bucket(), filesReq.Filename, filesReq)
	url, err := h.fileUsecase.UploadFiles(filesReq)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFilesErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		url,
	).Res()

}
