package users_repo

import (
	"todo-app/entity"
	"todo-app/pkg/errs"
)

type repoMock struct {
}

var (
	Add          func(user *entity.User) errs.Error
	FetchByEmail func(email string) (*entity.User, errs.Error)
	FetchById    func(userId uint) (*entity.User, errs.Error)
	Modify       func(userId uint, user *entity.User) errs.Error
)

func NewRepoMock() UsersRepo {
	return &repoMock{}
}

// Add implements UsersRepo.
func (rm *repoMock) Add(user *entity.User) errs.Error {
	return Add(user)
}

// FetchByEmail implements UsersRepo.
func (rm *repoMock) FetchByEmail(email string) (*entity.User, errs.Error) {
	return FetchByEmail(email)
}

// FetchById implements UsersRepo.
func (rm *repoMock) FetchById(userId uint) (*entity.User, errs.Error) {
	return FetchById(userId)
}

// Modify implements UsersRepo.
func (rm *repoMock) Modify(userId uint, user *entity.User) errs.Error {
	return Modify(userId, user)
}
