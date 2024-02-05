package usersHandlers

import (
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/entities"
	"github.com/NatthawutSK/NoTeams-Backend/modules/users"
	"github.com/NatthawutSK/NoTeams-Backend/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type userHandlerErrorCode string

const (
	signUpErr          userHandlerErrorCode = "user-001"
	signInErr          userHandlerErrorCode = "user-002"
	getUserProfileErr  userHandlerErrorCode = "user-003"
	signOutErr         userHandlerErrorCode = "user-004"
	refreshPassportErr userHandlerErrorCode = "user-005"
)

type IUsersHandler interface {
	SignIn(c *fiber.Ctx) error
	SignUp(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
}

type usersHandler struct {
	usersUsecases usersUsecases.IUserUsecase
	cfg           config.IConfig
}

func UsersHandler(usersUsecases usersUsecases.IUserUsecase, cfg config.IConfig) IUsersHandler {
	return &usersHandler{
		usersUsecases: usersUsecases,
		cfg:           cfg,
	}
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}

	result, err := h.usersUsecases.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *usersHandler) SignUp(c *fiber.Ctx) error {
	// Request Body parser
	req := new(users.UserRegisterReq)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpErr),
			err.Error(),
		).Res()
	}
	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpErr),
			"email is invalid",
		).Res()
	}

	// Insert user
	result, err := h.usersUsecases.InsertUser(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpErr),
				err.Error(),
			).Res()

		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpErr),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {

	userId := strings.Trim(c.Params("user_id"), " ")

	result, err := h.usersUsecases.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()

		}

	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	if err := h.usersUsecases.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecases.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}
