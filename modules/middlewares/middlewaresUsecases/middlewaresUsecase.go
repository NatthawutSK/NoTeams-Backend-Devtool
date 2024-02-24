package middlewaresUsecases

import "github.com/NatthawutSK/NoTeams-Backend/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	// IsMemberInTeam(userId, teamId string) bool
	// IsOwnerInTeam(userId, teamId string) bool
	IsAllowInviteMember(teamId string) bool
	IsAllowTask(teamId string) bool
	IsAllowFile(teamId string) bool
	AuthTeam(userId, teamId string) (bool, bool)
}

type middlewaresUsecase struct {
	middlewareRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewareRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewareRepository: middlewareRepository,
	}
}

func (u *middlewaresUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewareRepository.FindAccessToken(userId, accessToken)
}

// func (u *middlewaresUsecase) IsMemberInTeam(userId, teamId string) bool {
// 	return u.middlewareRepository.IsMemberInTeam(userId, teamId)
// }

// func (u *middlewaresUsecase) IsOwnerInTeam(userId, teamId string) bool {
// 	return u.middlewareRepository.IsOwnerInTeam(userId, teamId)
// }

func (u *middlewaresUsecase) IsAllowInviteMember(teamId string) bool {
	return u.middlewareRepository.IsAllowInviteMember(teamId)
}

func (u *middlewaresUsecase) AuthTeam(userId, teamId string) (bool, bool) {
	return u.middlewareRepository.AuthTeam(userId, teamId)
}

func (u *middlewaresUsecase) IsAllowTask(teamId string) bool {
	return u.middlewareRepository.IsAllowTask(teamId)
}

func (u *middlewaresUsecase) IsAllowFile(teamId string) bool {
	return u.middlewareRepository.IsAllowFile(teamId)
}
