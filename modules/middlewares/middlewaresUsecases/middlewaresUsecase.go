package middlewaresUsecases

import "github.com/NatthawutSK/NoTeams-Backend/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	CheckMemberInTeam(userId, teamId string) bool
	CheckOwnerInTeam(userId, teamId string) bool
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

func (u *middlewaresUsecase) CheckMemberInTeam(userId, teamId string) bool {
	return u.middlewareRepository.CheckMemberInTeam(userId, teamId)
}

func (u *middlewaresUsecase) CheckOwnerInTeam(userId, teamId string) bool {
	return u.middlewareRepository.CheckOwnerInTeam(userId, teamId)
}
