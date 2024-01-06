package todos_repo

import (
	"todo-app/entity"
	"todo-app/pkg/errs"
)

type repoMock struct {
}

var (
	Add    func(todo *entity.Todo) errs.Error
	Delete func(todoId uint) errs.Error
	Detail func(todoId uint) (*entity.Todo, errs.Error)
	Fetch  func(userId uint) ([]*entity.Todo, errs.Error)
	Modify func(todoId uint, todo *entity.Todo) errs.Error
)

func NewRepoMock() TodoRepo {
	return &repoMock{}
}

// Add implements TodoRepo.
func (rm *repoMock) Add(todo *entity.Todo) errs.Error {
	return Add(todo)
}

// Delete implements TodoRepo.
func (rm *repoMock) Delete(todoId uint) errs.Error {
	return Delete(todoId)
}

// Detail implements TodoRepo.
func (rm *repoMock) Detail(todoId uint) (*entity.Todo, errs.Error) {
	return Detail(todoId)
}

// Fetch implements TodoRepo.
func (rm *repoMock) Fetch(userId uint) ([]*entity.Todo, errs.Error) {
	return Fetch(userId)
}

// Modify implements TodoRepo.
func (rm *repoMock) Modify(todoId uint, todo *entity.Todo) errs.Error {
	return Modify(todoId, todo)
}
