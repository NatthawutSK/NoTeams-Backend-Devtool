package filesUsecase

import (
	"fmt"
	"math"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/utils"
)

type IFilesUsecase interface {
	UploadFiles(req []*multipart.FileHeader) error
}

type filesUsecase struct {
	cfg config.IConfig
}

func FilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}

func (u *filesUsecase) UploadFiles(filesReq []*multipart.FileHeader) error {
	filesUpload := make([]*files.FileReq, 0)

	// files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
		"pdf":  "pdf",
	}

	for _, file := range filesReq {
		// check file extension
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return fmt.Errorf("invalid file extension")
		}
		// check file size
		if file.Size > int64(u.cfg.App().FileLimit()) {
			return fmt.Errorf("file size must less than %d MB", int(math.Ceil(float64(u.cfg.App().FileLimit())/math.Pow(1024, 2))))
		}

		filename := utils.RandFileName(ext)
		filesUpload = append(filesUpload, &files.FileReq{
			FileName: filename,
			Files:    file,
		})
	}

	fmt.Println(filesUpload)

	return nil

}
