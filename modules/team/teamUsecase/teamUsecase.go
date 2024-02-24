package teamUsecase

import (
	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team"
	"github.com/NatthawutSK/NoTeams-Backend/modules/team/teamRepository"
)

type ITeamUsecase interface {
	CreateTeam(req *team.CreateTeamReq) (*team.CreateTeamRes, error)
	GetTeamById(teamId string) (*team.GetTeamByIdRes, error)
	JoinTeam(req *team.JoinTeamReq) (*team.JoinTeamRes, error)
	GetTeamByUserId(userId string) ([]*team.GetTeamByUserIdRes, error)
	InviteMember(team_id string, req *team.InviteMemberReq) error
}

type teamUsecase struct {
	teamRepo teamRepository.ITeamRepository
	cfg      config.IConfig
}

func TeamUsecase(teamRepo teamRepository.ITeamRepository, cfg config.IConfig) ITeamUsecase {
	return &teamUsecase{
		teamRepo: teamRepo,
		cfg:      cfg,
	}
}

func (u *teamUsecase) CreateTeam(req *team.CreateTeamReq) (*team.CreateTeamRes, error) {
	result, err := u.teamRepo.CreateTeam(req)
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
