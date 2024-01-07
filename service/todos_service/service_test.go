package todos_service_test

import (
	"testing"
	"time"
	"todo-app/dto"
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/repo/todos_repo"
	"todo-app/service/todos_service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var repoMock = todos_repo.NewRepoMock()
var service = todos_service.NewTodoService(repoMock)

var todoId = 1
var userId = 1

var add = &dto.AddTodo{
	Todos: "something message",
}

var modify = &dto.ModifyTodo{
	Todos:  "new todos",
	Status: true,
}

func TestAddTodoSuccess(t *testing.T) {
	todos_repo.Add = func(todo *entity.Todo) errs.Error {
		return nil
	}

	tr, err := service.Add(uint(todoId), add)

	assert.Nil(t, err)
	assert.NotNil(t, tr)
	assert.Equal(t, fiber.StatusCreated, tr.Status)
}

func TestAddTodoServerError(t *testing.T) {
	todos_repo.Add = func(todo *entity.Todo) errs.Error {
		return errs.NewInternalServerError("something went wrong")
	}

	tr, err := service.Add(uint(todoId), add)

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, err.Status())
}

func TestDeleteTodoDetailNotFound(t *testing.T) {
	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return nil, errs.NewNotFoundError("todo not found")
	}

	tr, err := service.Delete(uint(todoId))

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusNotFound, err.Status())
}

func TestDeleteTodoDetailServerError(t *testing.T) {
	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	tr, err := service.Delete(uint(todoId))

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, err.Status())
}

func TestDeleteTodoServerError(t *testing.T) {
	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return &entity.Todo{}, nil
	}

	todos_repo.Delete = func(todoId uint) errs.Error {
		return errs.NewInternalServerError("something went wrong")
	}

	tr, err := service.Delete(uint(todoId))

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, err.Status())
}

func TestDeleteTodoSuccess(t *testing.T) {
	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return &entity.Todo{}, nil
	}

	todos_repo.Delete = func(todoId uint) errs.Error {
		return nil
	}

	tr, err := service.Delete(uint(todoId))

	assert.Nil(t, err)
	assert.NotNil(t, tr)
	assert.Equal(t, fiber.StatusOK, tr.Status)
}

func TestDetailTodoServerError(t *testing.T) {
	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	tr, err := service.Detail(uint(todoId))

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, err.Status())
}

func TestDetailTodoNotFound(t *testing.T) {
	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return nil, errs.NewNotFoundError("todo not found")
	}

	tr, err := service.Detail(uint(todoId))

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusNotFound, err.Status())
}

func TestDetailTodoSuccess(t *testing.T) {
	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return &entity.Todo{}, nil
	}

	tr, err := service.Detail(uint(todoId))

	assert.Nil(t, err)
	assert.NotNil(t, tr)
	assert.Equal(t, fiber.StatusOK, tr.Status)
}

func TestFetchTodoServerError(t *testing.T) {
	todos_repo.Fetch = func(userId uint) ([]*entity.Todo, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	tr, err := service.Fetch(uint(userId))

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, err.Status())
}

func TestFetchTodoSuccess(t *testing.T) {
	todos_repo.Fetch = func(userId uint) ([]*entity.Todo, errs.Error) {
		return []*entity.Todo{
			{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}, nil
	}

	tr, err := service.Fetch(uint(userId))

	assert.Nil(t, err)
	assert.NotNil(t, tr)
	assert.Equal(t, fiber.StatusOK, tr.Status)
}

func TestModifyTodoDetailServerError(t *testing.T) {

	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	tr, err := service.Modify(uint(todoId), modify)

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, err.Status())
}

func TestModifyTodoDetailNotFound(t *testing.T) {

	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return nil, errs.NewNotFoundError("todo not found")
	}

	tr, err := service.Modify(uint(todoId), modify)

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusNotFound, err.Status())
}

func TestModifyTodoServerError(t *testing.T) {

	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return &entity.Todo{}, nil
	}

	todos_repo.Modify = func(todoId uint, todo *entity.Todo) errs.Error {
		return errs.NewInternalServerError("something went wrong")
	}

	tr, err := service.Modify(uint(todoId), modify)

	assert.Nil(t, tr)
	assert.NotNil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, err.Status())
}

func TestModifyTodoSuccess(t *testing.T) {

	todos_repo.Detail = func(todoId uint) (*entity.Todo, errs.Error) {
		return &entity.Todo{}, nil
	}

	todos_repo.Modify = func(todoId uint, todo *entity.Todo) errs.Error {
		return nil
	}

	tr, err := service.Modify(uint(todoId), modify)

	assert.Nil(t, err)
	assert.NotNil(t, tr)
	assert.Equal(t, fiber.StatusOK, tr.Status)
}
