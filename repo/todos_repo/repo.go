package todos_repo

import (
	"todo-app/entity"
	"todo-app/pkg/errs"
)

type TodoRepo interface {
	Add(todo *entity.Todo) errs.Error
	Fetch(userId uint) ([]*entity.Todo, errs.Error)
	Detail(todoId uint) (*entity.Todo, errs.Error)
	Modify(todoId uint, todo *entity.Todo) errs.Error
	Delete(todoId uint) errs.Error
}
