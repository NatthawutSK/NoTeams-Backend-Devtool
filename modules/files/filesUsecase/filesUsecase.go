package filesUsecase

import (
	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files/filesRepository"
)

type IFilesUsecase interface {
	GetFilesTeam(teamId string) (*files.GetFilesTeamRes, error)
	// UploadFiles(req []*multipart.FileHeader, isDownload bool, folder string) ([]*files.FileRes, error)
	// UploadFilesTeam(req []*multipart.FileHeader, isDownload bool, folder string) ([]*files.FileRes, error)
}

type filesUsecase struct {
	cfg       config.IConfig
	filesRepo filesRepository.IFilesRepository
}

func FilesUsecase(cfg config.IConfig, filesRepo filesRepository.IFilesRepository) IFilesUsecase {
	return &filesUsecase{
		cfg:       cfg,
		filesRepo: filesRepo,
	}
}

func (u *filesUsecase) GetFilesTeam(teamId string) (*files.GetFilesTeamRes, error) {

	files, err := u.filesRepo.GetFilesTeam(teamId)
	if err != nil {
		return nil, err
	}

	return files, nil
}
