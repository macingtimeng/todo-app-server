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

func TestFetchSuccess(t *testing.T) {

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	todos_service.Fetch = func(userId uint) (*dto.TodoResponse, errs.Error) {
		return &dto.TodoResponse{
			Status:  fiber.StatusOK,
			Message: "todos successfully fetched",
		}, nil
	}

	app.Get("/todos", auth_service.Authentication(), handler.Fetch)

	req := httptest.NewRequest(fiber.MethodGet, "/todos", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestFetchServerError(t *testing.T) {

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	todos_service.Fetch = func(userId uint) (*dto.TodoResponse, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	app.Get("/todos", auth_service.Authentication(), handler.Fetch)

	req := httptest.NewRequest(fiber.MethodGet, "/todos", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestDetailSuccess(t *testing.T) {

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	todos_service.Detail = func(todoId uint) (*dto.TodoResponse, errs.Error) {
		return &dto.TodoResponse{
			Status:  fiber.StatusOK,
			Message: "get todo with detail successfully fetched",
			Data: &dto.Todo{
				Id:    1,
				Todos: "todos",
			},
		}, nil
	}

	app.Get("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Detail)

	req := httptest.NewRequest(fiber.MethodGet, "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestDetailServerError(t *testing.T) {

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	todos_service.Detail = func(todoId uint) (*dto.TodoResponse, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	app.Get("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Detail)

	req := httptest.NewRequest(fiber.MethodGet, "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestDetailNotFound(t *testing.T) {

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	todos_service.Detail = func(todoId uint) (*dto.TodoResponse, errs.Error) {
		return nil, errs.NewNotFoundError("todo not found")
	}

	app.Get("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Detail)

	req := httptest.NewRequest(fiber.MethodGet, "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestDeleteSuccess(t *testing.T) {

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	todos_service.Delete = func(todoId uint) (*dto.TodoResponse, errs.Error) {
		return &dto.TodoResponse{
			Status:  fiber.StatusOK,
			Message: "todo successfully deleted",
		}, nil
	}

	app.Delete("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Delete)

	req := httptest.NewRequest(fiber.MethodDelete, "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestDeleteServerError(t *testing.T) {

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	todos_service.Delete = func(todoId uint) (*dto.TodoResponse, errs.Error) {
		return nil, errs.NewInternalServerError("something webt wrong")
	}

	app.Delete("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Delete)

	req := httptest.NewRequest(fiber.MethodDelete, "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestDeleteNotFound(t *testing.T) {

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	todos_service.Delete = func(todoId uint) (*dto.TodoResponse, errs.Error) {
		return nil, errs.NewNotFoundError("noy found")
	}

	app.Delete("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Delete)

	req := httptest.NewRequest(fiber.MethodDelete, "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestModifySuccess(t *testing.T) {

	b, _ := json.Marshal(modify)

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	todos_service.Modify = func(todoId uint, payload *dto.ModifyTodo) (*dto.TodoResponse, errs.Error) {
		return &dto.TodoResponse{
			Status:  fiber.StatusOK,
			Message: "todo successfully modifed",
		}, nil
	}

	app.Patch("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Modify)

	req := httptest.NewRequest(fiber.MethodPatch, "/todos/1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestModifyInvalidJSON(t *testing.T) {

	b, _ := json.Marshal(modify)

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	app.Patch("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Modify)

	req := httptest.NewRequest(fiber.MethodPatch, "/todos/1", bytes.NewReader(b))

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
}

func TestModifyBadRequest(t *testing.T) {

	b, _ := json.Marshal(&dto.ModifyTodo{})

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	app.Patch("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Modify)

	req := httptest.NewRequest(fiber.MethodPatch, "/todos/1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestModifyNotFound(t *testing.T) {

	b, _ := json.Marshal(modify)

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	todos_service.Modify = func(todoId uint, payload *dto.ModifyTodo) (*dto.TodoResponse, errs.Error) {
		return nil, errs.NewNotFoundError("todo not found")
	}

	app.Patch("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Modify)

	req := httptest.NewRequest(fiber.MethodPatch, "/todos/1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestModifyServerError(t *testing.T) {

	b, _ := json.Marshal(modify)

	auth_service.Authentication = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("user", user)
			return c.Next()
		}
	}

	auth_service.Authorization = func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	todos_service.Modify = func(todoId uint, payload *dto.ModifyTodo) (*dto.TodoResponse, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	app.Patch("/todos/:todoId", auth_service.Authentication(), auth_service.Authorization(), handler.Modify)

	req := httptest.NewRequest(fiber.MethodPatch, "/todos/1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, 1)

	assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}
