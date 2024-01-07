package todos_handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"todo-app/dto"
	"todo-app/entity"
	"todo-app/handler/todos_handler"
	"todo-app/pkg/errs"
	"todo-app/service/auth_service"
	"todo-app/service/todos_service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var serviceMock = todos_service.NewServiceMock()
var handler = todos_handler.NewTodoHandler(serviceMock)

var app = fiber.New()

var userId = 1
var todoId = 1

var add = &dto.AddTodo{
	Todos: "todos",
}

var modify = &dto.ModifyTodo{
	Todos:  "new todos",
	Status: true,
}

var user = entity.User{
	Model: gorm.Model{
		ID: 1,
	},
	Name:  "jihan",
	Email: "jihan@weeekly.com",
}

func TestAddSuccess(t *testing.T) {
	b, _ := json.Marshal(add)

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	todos_service.Add = func(userId uint, payload *dto.AddTodo) (*dto.TodoResponse, errs.Error) {
		return &dto.TodoResponse{
			Status:  fiber.StatusCreated,
			Message: "todo successfully added",
		}, nil
	}

	app.Post("/todos", auth_service.Authentication(), handler.Add)

	req := httptest.NewRequest(fiber.MethodPost, "/todos", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestAddServerError(t *testing.T) {
	b, _ := json.Marshal(add)

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	todos_service.Add = func(userId uint, payload *dto.AddTodo) (*dto.TodoResponse, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	app.Post("/todos", auth_service.Authentication(), handler.Add)

	req := httptest.NewRequest(fiber.MethodPost, "/todos", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestAddInvalidJSON(t *testing.T) {
	b, _ := json.Marshal(add)

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	app.Post("/todos", auth_service.Authentication(), handler.Add)

	req := httptest.NewRequest(fiber.MethodPost, "/todos", bytes.NewReader(b))

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
}

func TestAddBadRequest(t *testing.T) {
	b, _ := json.Marshal(&dto.AddTodo{})

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	app.Post("/todos", auth_service.Authentication(), handler.Add)

	req := httptest.NewRequest(fiber.MethodPost, "/todos", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}
