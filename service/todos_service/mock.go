package todos_service

import (
	"todo-app/dto"
	"todo-app/pkg/errs"
)

type serviceMock struct {
}

var (
	Add    func(userId uint, payload *dto.AddTodo) (*dto.TodoResponse, errs.Error)
	Delete func(todoId uint) (*dto.TodoResponse, errs.Error)
	Detail func(todoId uint) (*dto.TodoResponse, errs.Error)
	Fetch  func(userId uint) (*dto.TodoResponse, errs.Error)
	Modify func(todoId uint, payload *dto.ModifyTodo) (*dto.TodoResponse, errs.Error)
)

func NewServiceMock() TodoService {
	return &serviceMock{}
}

// Add implements TodoService.
func (sm *serviceMock) Add(userId uint, payload *dto.AddTodo) (*dto.TodoResponse, errs.Error) {
	return Add(userId, payload)
}

// Delete implements TodoService.
func (sm *serviceMock) Delete(todoId uint) (*dto.TodoResponse, errs.Error) {
	return Delete(todoId)
}

// Detail implements TodoService.
func (sm *serviceMock) Detail(todoId uint) (*dto.TodoResponse, errs.Error) {
	return Detail(todoId)
}

// Fetch implements TodoService.
func (sm *serviceMock) Fetch(userId uint) (*dto.TodoResponse, errs.Error) {
	return Fetch(userId)
}

// Modify implements TodoService.
func (sm *serviceMock) Modify(todoId uint, payload *dto.ModifyTodo) (*dto.TodoResponse, errs.Error) {
	return Modify(todoId, payload)
}
