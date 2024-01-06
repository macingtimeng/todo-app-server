package todos_service

import (
	"todo-app/dto"
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/repo/todos_repo"

	"github.com/gofiber/fiber/v2"
)

type todoService struct {
	tr todos_repo.TodoRepo
}

type TodoService interface {
	Add(userId uint, payload *dto.AddTodo) (*dto.TodoResponse, errs.Error)
	Delete(todoId uint) (*dto.TodoResponse, errs.Error)
	Detail(todoId uint) (*dto.TodoResponse, errs.Error)
	Fetch(userId uint) (*dto.TodoResponse, errs.Error)
	Modify(todoId uint, payload *dto.ModifyTodo) (*dto.TodoResponse, errs.Error)
}

func NewTodoService(todoRepo todos_repo.TodoRepo) TodoService {
	return &todoService{tr: todoRepo}
}

// Add implements TodoService.
func (ts *todoService) Add(userId uint, payload *dto.AddTodo) (*dto.TodoResponse, errs.Error) {

	err := ts.tr.Add(&entity.Todo{
		Todos:  payload.Todos,
		UserID: userId,
	})

	if err != nil {
		return nil, err
	}

	return &dto.TodoResponse{
		Status:  fiber.StatusCreated,
		Message: "todo successfully added",
		Data:    nil,
	}, nil
}

// Delete implements TodoService.
func (ts *todoService) Delete(todoId uint) (*dto.TodoResponse, errs.Error) {

	_, err := ts.tr.Detail(todoId)

	if err != nil {
		return nil, err
	}

	err = ts.tr.Delete(todoId)

	if err != nil {
		return nil, err
	}

	return &dto.TodoResponse{
		Status:  fiber.StatusOK,
		Message: "todo successfully deleted",
		Data:    nil,
	}, nil
}

// Detail implements TodoService.
func (ts *todoService) Detail(todoId uint) (*dto.TodoResponse, errs.Error) {

	todo, err := ts.tr.Detail(todoId)

	if err != nil {
		return nil, err
	}

	return &dto.TodoResponse{
		Status:  fiber.StatusOK,
		Message: "get todo with detail successfully fetched",
		Data: &dto.Todo{
			Id:        todo.ID,
			Todos:     todo.Todos,
			Status:    todo.Status,
			CreatedAt: todo.CreatedAt,
			UpdatedAt: todo.UpdatedAt,
		},
	}, nil
}

// Fetch implements TodoService.
func (ts *todoService) Fetch(userId uint) (*dto.TodoResponse, errs.Error) {

	t, err := ts.tr.Fetch(userId)

	if err != nil {
		return nil, err
	}

	todos := []*dto.Todo{}

	for _, eachTodo := range t {
		todo := &dto.Todo{
			Id:        eachTodo.ID,
			Todos:     eachTodo.Todos,
			Status:    eachTodo.Status,
			CreatedAt: eachTodo.CreatedAt,
			UpdatedAt: eachTodo.UpdatedAt,
		}

		todos = append(todos, todo)
	}

	return &dto.TodoResponse{
		Status:  fiber.StatusOK,
		Message: "todos successfully fetched",
		Data:    todos,
	}, nil
}

// Modify implements TodoService.
func (ts *todoService) Modify(todoId uint, payload *dto.ModifyTodo) (*dto.TodoResponse, errs.Error) {

	_, err := ts.tr.Detail(todoId)

	if err != nil {
		return nil, err
	}

	err = ts.tr.Modify(todoId, payload.ModifyTodoToEntity())

	if err != nil {
		return nil, err
	}

	return &dto.TodoResponse{
		Status:  fiber.StatusOK,
		Message: "todo successfully modified",
		Data:    nil,
	}, nil
}
