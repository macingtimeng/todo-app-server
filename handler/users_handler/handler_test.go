package users_handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"todo-app/dto"
	"todo-app/entity"
	"todo-app/handler/users_handler"
	"todo-app/pkg/errs"
	"todo-app/service/auth_service"
	"todo-app/service/users_service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var serviceMock = users_service.NewServiceMock()
var authMock = auth_service.NewAuthMock()

var app = fiber.New()
var userHandler = users_handler.NewUserHandler(serviceMock)

var register = &dto.Register{
	Name:     "Jihan",
	Email:    "jihan@weeekly.com",
	Password: "secret",
}

var login = &dto.Login{
	Email:    "jihan@weeekly.com",
	Password: "secret",
}

var modify = &dto.Modify{
	Name:  "Jihan Weeekly",
	Email: "jihan@weeekly.com",
}

var user = entity.User{
	Model: gorm.Model{
		ID: 1,
	},
	Name:  "Jihan",
	Email: "jihan@weeekly.com",
}

func TestRegisterSuccess(t *testing.T) {
	b, _ := json.Marshal(register)

	users_service.Register = func(payload *dto.Register) (*dto.UserResponse, errs.Error) {
		return &dto.UserResponse{
			Status:  fiber.StatusCreated,
			Message: "user successfully created",
			Data:    nil,
		}, nil
	}

	app.Post("/users/register", userHandler.Register)

	req := httptest.NewRequest(fiber.MethodPost, "/users/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestRegisterServerError(t *testing.T) {
	b, _ := json.Marshal(register)

	users_service.Register = func(payload *dto.Register) (*dto.UserResponse, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	app.Post("/users/register", userHandler.Register)

	req := httptest.NewRequest(fiber.MethodPost, "/users/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestRegisterConflictError(t *testing.T) {
	b, _ := json.Marshal(register)

	users_service.Register = func(payload *dto.Register) (*dto.UserResponse, errs.Error) {
		return nil, errs.NewConflictError("email has been used")
	}

	app.Post("/users/register", userHandler.Register)

	req := httptest.NewRequest(fiber.MethodPost, "/users/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusConflict, res.StatusCode)
}

func TestRegisterInvalidJSON(t *testing.T) {
	b, _ := json.Marshal(&dto.Register{})

	app.Post("/users/register", userHandler.Register)

	req := httptest.NewRequest(fiber.MethodPost, "/users/register", bytes.NewReader(b))

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
}

func TestRegisterBadRequest(t *testing.T) {
	b, _ := json.Marshal(&dto.Register{})

	app.Post("/users/register", userHandler.Register)

	req := httptest.NewRequest(fiber.MethodPost, "/users/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestLoginSuccess(t *testing.T) {
	b, _ := json.Marshal(login)

	users_service.Login = func(payload *dto.Login) (*dto.UserResponse, errs.Error) {
		return &dto.UserResponse{
			Status:  fiber.StatusOK,
			Message: "user successfully loged in",
			Data: dto.Token{
				TokenString: user.GenerateToken(),
			},
		}, nil
	}

	app.Post("/users/login", userHandler.Login)

	req := httptest.NewRequest(fiber.MethodPost, "/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestLoginInvalidUser(t *testing.T) {
	b, _ := json.Marshal(login)

	users_service.Login = func(payload *dto.Login) (*dto.UserResponse, errs.Error) {
		return nil, errs.NewUnauthenticatedError("invalid email or password")
	}

	app.Post("/users/login", userHandler.Login)

	req := httptest.NewRequest(fiber.MethodPost, "/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestLoginServerError(t *testing.T) {
	b, _ := json.Marshal(login)

	users_service.Login = func(payload *dto.Login) (*dto.UserResponse, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	app.Post("/users/login", userHandler.Login)

	req := httptest.NewRequest(fiber.MethodPost, "/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestLoginBadRequest(t *testing.T) {
	b, _ := json.Marshal(&dto.Login{})

	app.Post("/users/login", userHandler.Login)

	req := httptest.NewRequest(fiber.MethodPost, "/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestLoginInvalidJSON(t *testing.T) {
	b, _ := json.Marshal(&dto.Login{})

	app.Post("/users/login", userHandler.Login)

	req := httptest.NewRequest(fiber.MethodPost, "/users/login", bytes.NewReader(b))

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
}

func TestModifySuccess(t *testing.T) {
	b, _ := json.Marshal(modify)

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	users_service.Modify = func(userId uint, payload *dto.Modify) (*dto.UserResponse, errs.Error) {
		return &dto.UserResponse{
			Status:  fiber.StatusOK,
			Message: "user successfully modified",
			Data:    nil,
		}, nil
	}

	app.Patch("/users/modify", authMock.Authentication(), userHandler.Modify)

	req := httptest.NewRequest(fiber.MethodPatch, "/users/modify", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestModifyServerError(t *testing.T) {
	b, _ := json.Marshal(modify)

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	users_service.Modify = func(userId uint, payload *dto.Modify) (*dto.UserResponse, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	app.Patch("/users/modify", authMock.Authentication(), userHandler.Modify)

	req := httptest.NewRequest(fiber.MethodPatch, "/users/modify", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestModifyBadRequest(t *testing.T) {
	b, _ := json.Marshal(&dto.Modify{})

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	app.Patch("/users/modify", authMock.Authentication(), userHandler.Modify)

	req := httptest.NewRequest(fiber.MethodPatch, "/users/modify", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestModifyInvalidJSON(t *testing.T) {
	b, _ := json.Marshal(&dto.Modify{})

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	app.Patch("/users/modify", authMock.Authentication(), userHandler.Modify)

	req := httptest.NewRequest(fiber.MethodPatch, "/users/modify", bytes.NewReader(b))

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
}

func TestProfileSuccess(t *testing.T) {
	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	users_service.Profile = func(userId uint) (*dto.UserResponse, errs.Error) {
		return &dto.UserResponse{
			Status:  fiber.StatusOK,
			Message: "user successfully fetched",
			Data:    &dto.Profile{},
		}, nil
	}

	app.Get("/users/profile", authMock.Authentication(), userHandler.Profile)

	req := httptest.NewRequest(fiber.MethodGet, "/users/profile", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestProfileNotFound(t *testing.T) {
	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	users_service.Profile = func(userId uint) (*dto.UserResponse, errs.Error) {
		return nil, errs.NewNotFoundError("user not found")
	}

	app.Get("/users/profile", authMock.Authentication(), userHandler.Profile)

	req := httptest.NewRequest(fiber.MethodGet, "/users/profile", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestProfileServerError(t *testing.T) {
	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	users_service.Profile = func(userId uint) (*dto.UserResponse, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	app.Get("/users/profile", authMock.Authentication(), userHandler.Profile)

	req := httptest.NewRequest(fiber.MethodGet, "/users/profile", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}
