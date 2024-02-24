package middlewaresHandlers

import (
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/entities"
	"github.com/NatthawutSK/NoTeams-Backend/modules/middlewares/middlewaresUsecases"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middleware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middleware-002"
	paramsCheckErr middlewareHandlersErrCode = "middleware-003"
	authTeam       middlewareHandlersErrCode = "middleware-004"
	IsAllowInvite  middlewareHandlersErrCode = "middleware-005"
	IsAllowTask    middlewareHandlersErrCode = "middleware-006"
	IsAllowFile    middlewareHandlersErrCode = "middleware-007"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	// IsMemberTeam() fiber.Handler
	// IsOwnerTeam() fiber.Handler
	IsAllowInvite() fiber.Handler
	IsAllowTask() fiber.Handler
	IsAllowFile() fiber.Handler
	AuthTeam() fiber.Handler
}

type middlewaresHandler struct {
	cfg                config.IConfig
	middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                cfg,
		middlewaresUsecase: middlewaresUsecase,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

// ตรวจสอบว่ามี router ที่เรียกหรือไม่
func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

// แกะ token และตรวจสอบว่า Login อยู่หรือไม่
func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		// check token in db
		check := h.middlewaresUsecase.FindAccessToken(claims.Id, token)
		if !check {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"You Are Not Logged In",
			).Res()
		}

		// set UserId and roleId to locals
		c.Locals("userId", claims.Id)
		return c.Next()
	}
}

// ป้องกันการเข้าถึงข้อมูลของคนอื่น ต้องมาคู่กับ JwtAuth
func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"no permission to access",
			).Res()
		}
		return c.Next()
	}

}

// // ตรวจสอบว่าเป็นสมาชิกของทีมหรือไม่ ต้องมาคู่กับ JwtAuth
// func (h *middlewaresHandler) IsMemberTeam() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		userId := c.Locals("userId").(string)
// 		teamId := strings.TrimSpace(c.Params("team_id"))
// 		check := h.middlewaresUsecase.IsMemberInTeam(userId, teamId)
// 		if !check {
// 			return entities.NewResponse(c).Error(
// 				fiber.ErrUnauthorized.Code,
// 				string(authTeam),
// 				"no permission to access team",
// 			).Res()
// 		}
// 		return c.Next()
// 	}
// }

// // ตรวจสอบว่าเป็นเจ้าของทีมหรือไม่ ต้องมาคู่กับ JwtAuth
// func (h *middlewaresHandler) IsOwnerTeam() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		userId := c.Locals("userId").(string)
// 		teamId := strings.TrimSpace(c.Params("team_id"))
// 		check := h.middlewaresUsecase.IsOwnerInTeam(userId, teamId)
// 		if !check {
// 			return entities.NewResponse(c).Error(
// 				fiber.ErrUnauthorized.Code,
// 				string(authTeam),
// 				"only owner have permission",
// 			).Res()
// 		}
// 		// c.Locals("role", "OWNER")
// 		return c.Next()
// 	}
// }

// ตรวจสอบว่าเป็นสมาชิกหรือเจ้าของทีมหรือไม่ ต้องมาคู่กับ JwtAuth
func (h *middlewaresHandler) AuthTeam() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(string)
		teamId := strings.TrimSpace(c.Params("team_id"))
		isMember, isOwner := h.middlewaresUsecase.AuthTeam(userId, teamId)
		if !isMember {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authTeam),
				"no permission to access team",
			).Res()
		}
		c.Locals("role", "MEMBER")
		if isOwner {
			c.Locals("role", "OWNER")
		}
		return c.Next()
	}
}

// ตรวจสอบว่าสามารถเชิญคนเข้าทีมได้หรือไม่ ต้องมาคู่กับ JwtAuth และ AuthTeam
func (h *middlewaresHandler) IsAllowInvite() fiber.Handler {
	return func(c *fiber.Ctx) error {
		teamId := strings.TrimSpace(c.Params("team_id"))
		role := c.Locals("role").(string)
		if role == "OWNER" {
			return c.Next()
		}
		check := h.middlewaresUsecase.IsAllowInviteMember(teamId)
		if !check {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(IsAllowInvite),
				"no permission to invite member",
			).Res()
		}

		return c.Next()
	}
}

// ตรวจสอบว่าสามารถสร้างหรือแก้ไข task ได้หรือไม่ ต้องมาคู่กับ JwtAuth และ AuthTeam
func (h *middlewaresHandler) IsAllowTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		teamId := strings.TrimSpace(c.Params("team_id"))
		role := c.Locals("role").(string)
		if role == "OWNER" {
			return c.Next()
		}
		check := h.middlewaresUsecase.IsAllowTask(teamId)
		if !check {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(IsAllowTask),
				"no permission to manage task",
			).Res()
		}

		return c.Next()
	}
}

// ตรวจสอบว่าสามารถอัพโหลดไฟล์ได้หรือไม่ ต้องมาคู่กับ JwtAuth และ AuthTeam
func (h *middlewaresHandler) IsAllowFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		teamId := strings.TrimSpace(c.Params("team_id"))
		role := c.Locals("role").(string)
		if role == "OWNER" {
			return c.Next()
		}
		check := h.middlewaresUsecase.IsAllowFile(teamId)
		if !check {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(IsAllowFile),
				"no permission to manage file",
			).Res()
		}

		return c.Next()
	}
}
