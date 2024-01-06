package users_repo

import (
	"todo-app/entity"
	"todo-app/pkg/errs"
)

type UsersRepo interface {
	Add(user *entity.User) errs.Error
	FetchById(userId uint) (*entity.User, errs.Error)
	FetchByEmail(email string) (*entity.User, errs.Error)
	Modify(userId uint, user *entity.User) errs.Error
}
