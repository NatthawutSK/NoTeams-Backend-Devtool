package filesUsecase

import (
	"mime/multipart"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files/filesRepository"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/utils"
)

type IFilesUsecase interface {
	GetFilesTeam(teamId string) (*files.GetFilesTeamRes, error)
	UploadFilesTeam(userId string, teamId string, filesReq []*multipart.FileHeader) ([]*files.FileTeamByIdRes, error)
	DeleteFilesTeam(req *files.DeleteFilesTeamReq) error
}

type filesUsecase struct {
	cfg       config.IConfig
	filesRepo filesRepository.IFilesRepository
	upload    utils.IUpload
}

func FilesUsecase(cfg config.IConfig, filesRepo filesRepository.IFilesRepository) IFilesUsecase {
	return &filesUsecase{
		cfg:       cfg,
		filesRepo: filesRepo,
		upload:    utils.Upload(cfg),
	}
}

func (u *filesUsecase) GetFilesTeam(teamId string) (*files.GetFilesTeamRes, error) {

	files, err := u.filesRepo.GetFilesTeam(teamId)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (u *filesUsecase) UploadFilesTeam(userId string, teamId string, filesReq []*multipart.FileHeader) ([]*files.FileTeamByIdRes, error) {

	folder := teamId + "/files"

	files, err := u.upload.UploadFiles(filesReq, true, folder)
	if err != nil {
		return nil, err
	}

	filesRes, err := u.filesRepo.UploadFilesTeam(userId, teamId, files)
	if err != nil {
		return nil, err
	}

	return filesRes, nil
}

func (u *filesUsecase) DeleteFilesTeam(req *files.DeleteFilesTeamReq) error {

	err := u.filesRepo.DeleteFilesTeam(req)
	if err != nil {
		return err
	}

	return nil
}
