package teamUsecase

import (
	"fmt"
	"mime/multipart"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team/teamRepository"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/utils"
)

type ITeamUsecase interface {
	CreateTeam(userId string, req *team.CreateTeamReq) (*team.CreateTeamRes, error)
	GetTeamById(teamId string) (*team.GetTeamByIdRes, error)
	JoinTeam(req *team.JoinTeamReq) (*team.JoinTeamRes, error)
	GetTeamByUserId(userId string) ([]*team.GetTeamByUserIdRes, error)
	InviteMember(team_id string, req *team.InviteMemberReq) error
	GetMemberTeam(teamId string) ([]*team.GetMemberTeamRes, error)
	DeleteMember(memberId string) error
	GetAboutTeam(teamId string) (*team.GetAboutTeamRes, error)
	GetSettingTeam(teamId string) (*team.GetSettingTeamRes, error)
	UpdateProfileTeam(userId string, req *team.UpdateTeamReq, posterFile []*multipart.FileHeader) error
	UpdatePermission(teamId string, req *team.UpdatePermissionReq) error
	UpdateCodeTeam(teamId string, req *team.UpdateCodeTeamReq) error
	DeleteTeam(teamId string) error
	ExitTeam(userId, teamId string) error
}

type teamUsecase struct {
	teamRepo teamRepository.ITeamRepository
	cfg      config.IConfig
	upload   utils.IUpload
}

func TeamUsecase(teamRepo teamRepository.ITeamRepository, cfg config.IConfig) ITeamUsecase {
	return &teamUsecase{
		teamRepo: teamRepo,
		cfg:      cfg,
		upload:   utils.Upload(cfg),
	}
}

func (u *teamUsecase) CreateTeam(userId string, req *team.CreateTeamReq) (*team.CreateTeamRes, error) {
	result, err := u.teamRepo.CreateTeam(userId, req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *teamUsecase) GetTeamById(teamId string) (*team.GetTeamByIdRes, error) {
	result, err := u.teamRepo.GetTeamById(teamId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *teamUsecase) JoinTeam(req *team.JoinTeamReq) (*team.JoinTeamRes, error) {
	result, err := u.teamRepo.JoinTeam(req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *teamUsecase) GetTeamByUserId(userId string) ([]*team.GetTeamByUserIdRes, error) {
	result, err := u.teamRepo.GetTeamByUserId(userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *teamUsecase) InviteMember(team_id string, req *team.InviteMemberReq) error {
	err := u.teamRepo.InviteMember(team_id, req)
	if err != nil {
		return err
	}
	return nil
}

func (u *teamUsecase) GetMemberTeam(teamId string) ([]*team.GetMemberTeamRes, error) {
	result, err := u.teamRepo.GetMemberTeam(teamId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *teamUsecase) DeleteMember(memberId string) error {
	err := u.teamRepo.DeleteMember(memberId)
	if err != nil {
		return err
	}
	return nil
}

func (u *teamUsecase) GetAboutTeam(teamId string) (*team.GetAboutTeamRes, error) {
	result, err := u.teamRepo.GetAboutTeam(teamId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *teamUsecase) GetSettingTeam(teamId string) (*team.GetSettingTeamRes, error) {
	result, err := u.teamRepo.GetSettingTeam(teamId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *teamUsecase) UpdateProfileTeam(teamId string, req *team.UpdateTeamReq, posterFile []*multipart.FileHeader) error {
	//check avatarFile

	if len(posterFile) > 0 {
		//check len of posterFile must be 1
		if len(posterFile) > 1 {
			return fmt.Errorf("poster file must be 1")
		}

		//upload poster
		url, err := u.upload.UploadFiles(posterFile, false, teamId)
		if err != nil {
			return err
		}

		req.TeamPoster = url[0].Url
	}

	//update profile team
	if err := u.teamRepo.UpdateTeam(teamId, req); err != nil {
		return err
	}

	return nil
}

func (u *teamUsecase) UpdatePermission(teamId string, req *team.UpdatePermissionReq) error {
	// check permissionType must be task | file | invite use map
	permissionTypeMap := map[string]string{
		"task":   "allow_task",
		"file":   "allow_file",
		"invite": "allow_invite",
	}

	//check permissionType
	if _, ok := permissionTypeMap[req.PermissionType]; !ok {
		return fmt.Errorf("invalid permission type")
	}

	//update type permission
	reqUpdate := &team.UpdatePermissionReq{
		PermissionType: permissionTypeMap[req.PermissionType],
		Value:          req.Value,
	}

	//update permission
	err := u.teamRepo.UpdatePermission(teamId, reqUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (u *teamUsecase) UpdateCodeTeam(teamId string, req *team.UpdateCodeTeamReq) error {
	//update code team
	err := u.teamRepo.UpdateCodeTeam(teamId, req)
	if err != nil {
		return err
	}
	return nil
}

func (u *teamUsecase) DeleteTeam(teamId string) error {
	//delete team
	err := u.teamRepo.DeleteTeam(teamId)
	if err != nil {
		return err
	}
	return nil
}

func (u *teamUsecase) ExitTeam(userId, teamId string) error {
	//exit team
	err := u.teamRepo.ExitTeam(userId, teamId)
	if err != nil {
		return err
	}
	return nil
}
