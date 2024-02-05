package middlewaresUsecases

import "github.com/NatthawutSK/NoTeams-Backend/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
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
