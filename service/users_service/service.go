package users_service

import (
	"net/http"
	"todo-app/dto"
	"todo-app/pkg/errs"
	"todo-app/repo/users_repo"
)

type userService struct {
	ur users_repo.UsersRepo
}

type UserService interface {
	Register(payload *dto.Register) (*dto.UserResponse, errs.Error)
	Login(payload *dto.Login) (*dto.UserResponse, errs.Error)
	Profile(userId uint) (*dto.UserResponse, errs.Error)
	Modify(userId uint, payload *dto.Modify) (*dto.UserResponse, errs.Error)
}

func NewUserService(userRepo users_repo.UsersRepo) UserService {
	return &userService{ur: userRepo}
}

// Login implements UserService.
func (us *userService) Login(payload *dto.Login) (*dto.UserResponse, errs.Error) {

	u, err := us.ur.FetchByEmail(payload.Email)

	if err != nil {
		if err.Status() == http.StatusNotFound {
			return nil, errs.NewUnauthenticatedError("invalid user email or password")
		}
		return nil, err
	}

	isValidPassword := u.CompareHashPassword(payload.Password)

	if !isValidPassword {
		return nil, errs.NewUnauthenticatedError("invalid user email or password")
	}

	return &dto.UserResponse{
		Status:  http.StatusOK,
		Message: "user successfully loged in",
		Data: dto.Token{
			TokenString: u.GenerateToken(),
		},
	}, nil
}

// Modify implements UserService.
func (us *userService) Modify(userId uint, payload *dto.Modify) (*dto.UserResponse, errs.Error) {

	err := us.ur.Modify(userId, payload.ModifyToEntity())

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Status:  http.StatusOK,
		Message: "user successfully modified",
		Data:    nil,
	}, nil
}

// Profile implements UserService.
func (us *userService) Profile(userId uint) (*dto.UserResponse, errs.Error) {

	u, err := us.ur.FetchById(userId)

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Status:  http.StatusOK,
		Message: "user successfully fetched",
		Data: dto.Profile{
			Id:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		},
	}, nil
}

// Register implements UserService.
func (us *userService) Register(payload *dto.Register) (*dto.UserResponse, errs.Error) {

	user := payload.RegisterToEntity()
	user.HashPassword()

	err := us.ur.Add(user)

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Status:  http.StatusCreated,
		Message: "user successfully created",
		Data:    nil,
	}, nil
}
