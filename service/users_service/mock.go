package users_service

import (
	"todo-app/dto"
	"todo-app/pkg/errs"
)

type serviceMock struct {
}

var (
	Login    func(payload *dto.Login) (*dto.UserResponse, errs.Error)
	Modify   func(userId uint, payload *dto.Modify) (*dto.UserResponse, errs.Error)
	Profile  func(userId uint) (*dto.UserResponse, errs.Error)
	Register func(payload *dto.Register) (*dto.UserResponse, errs.Error)
)

func NewServiceMock() UserService {
	return &serviceMock{}
}

// Login implements UserService.
func (sm *serviceMock) Login(payload *dto.Login) (*dto.UserResponse, errs.Error) {
	return Login(payload)
}

// Modify implements UserService.
func (sm *serviceMock) Modify(userId uint, payload *dto.Modify) (*dto.UserResponse, errs.Error) {
	return Modify(userId, payload)
}

// Profile implements UserService.
func (sm *serviceMock) Profile(userId uint) (*dto.UserResponse, errs.Error) {
	return Profile(userId)
}

// Register implements UserService.
func (sm *serviceMock) Register(payload *dto.Register) (*dto.UserResponse, errs.Error) {
	return Register(payload)
}
