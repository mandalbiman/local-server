package handlers

import (
	"bats.com/local-server/api/models"
	"bats.com/local-server/api/svc"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Auth interface {
	Login(c *fiber.Ctx) error
}

func NewAuthHandler(opts *AuthHandlerOpts) Auth {
	return &AuthHandlerImpl{
		validate: opts.Validate,
		authSvc:  opts.AuthSvc,
	}
}

func SetUpAuthRoutes(v1 fiber.Router) {
	validate := validator.New()
	authSvc := svc.NewAuth()
	uh := NewAuthHandler(&AuthHandlerOpts{Validate: validate, AuthSvc: authSvc})
	v1.Post("/login", uh.Login)
}

type AuthHandlerOpts struct {
	Validate *validator.Validate
	AuthSvc  svc.Auth
}

type AuthHandlerImpl struct {
	validate *validator.Validate
	authSvc  svc.Auth
}

func (a *AuthHandlerImpl) Login(c *fiber.Ctx) error {
	req := &models.LoginRequest{}
	if err := parseAndValidateRequest(c, req, a.validate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	res, err := a.authSvc.Login(req)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.JSON(res)
}
