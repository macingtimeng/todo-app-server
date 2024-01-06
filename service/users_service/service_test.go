package users_service_test

import (
	"net/http"
	"testing"
	"todo-app/dto"
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/repo/users_repo"
	"todo-app/service/users_service"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var repoMock = users_repo.NewRepoMock()
var service = users_service.NewUserService(repoMock)

var register = &dto.Register{
	Name:     "Jihan",
	Email:    "jihan@weeekly.com",
	Password: "secret",
}

var login = &dto.Login{
	Email:    "jihan@weeekly.com",
	Password: "secret",
}

var modify = &dto.Modify{
	Name:  "Jihan Weeekly",
	Email: "jihan@weeekly.com",
}

var userId = 1

func TestAddUserSuccess(t *testing.T) {
	users_repo.Add = func(user *entity.User) errs.Error {
		return nil
	}

	ur, err := service.Register(register)

	assert.Nil(t, err)
	assert.NotNil(t, ur)
	assert.Equal(t, http.StatusCreated, ur.Status)
}

func TestAddUserServerError(t *testing.T) {
	users_repo.Add = func(user *entity.User) errs.Error {
		return errs.NewInternalServerError("something went wrong")
	}

	ur, err := service.Register(register)

	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestAddUserConflict(t *testing.T) {
	users_repo.Add = func(user *entity.User) errs.Error {
		return errs.NewConflictError("email has been used")
	}

	ur, err := service.Register(register)

	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Status())
}

func TestProfileSuccess(t *testing.T) {
	users_repo.FetchById = func(userId uint) (*entity.User, errs.Error) {
		return &entity.User{}, nil
	}

	ur, err := service.Profile(uint(userId))

	assert.Nil(t, err)
	assert.NotNil(t, ur)
	assert.Equal(t, http.StatusOK, ur.Status)
}

func TestProfileNotFound(t *testing.T) {
	users_repo.FetchById = func(userId uint) (*entity.User, errs.Error) {
		return nil, errs.NewNotFoundError("user not found")
	}

	ur, err := service.Profile(uint(userId))

	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestProfileServerError(t *testing.T) {
	users_repo.FetchById = func(userId uint) (*entity.User, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	ur, err := service.Profile(uint(userId))

	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestModifySuccess(t *testing.T) {
	users_repo.Modify = func(userId uint, user *entity.User) errs.Error {
		return nil
	}

	ur, err := service.Modify(uint(userId), modify)

	assert.Nil(t, err)
	assert.NotNil(t, ur)
	assert.Equal(t, http.StatusOK, ur.Status)
}

func TestModifyServerError(t *testing.T) {
	users_repo.Modify = func(userId uint, user *entity.User) errs.Error {
		return errs.NewInternalServerError("something went wrong")
	}

	ur, err := service.Modify(uint(userId), modify)

	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestLoginNotFound(t *testing.T) {
	users_repo.FetchByEmail = func(email string) (*entity.User, errs.Error) {
		return nil, errs.NewNotFoundError("user not found")
	}

	ur, err := service.Login(login)

	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusUnauthorized, err.Status())
}

func TestLoginServerError(t *testing.T) {
	users_repo.FetchByEmail = func(email string) (*entity.User, errs.Error) {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	ur, err := service.Login(login)

	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestLoginInvalidPassword(t *testing.T) {
	users_repo.FetchByEmail = func(email string) (*entity.User, errs.Error) {
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)

		return &entity.User{
			Password: string(hashPassword),
		}, nil
	}

	ur, err := service.Login(&dto.Login{Email: "", Password: ""})

	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusUnauthorized, err.Status())
}

func TestLoginSuccess(t *testing.T) {
	users_repo.FetchByEmail = func(email string) (*entity.User, errs.Error) {
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)

		return &entity.User{
			Password: string(hashPassword),
		}, nil
	}

	ur, err := service.Login(login)

	assert.Nil(t, err)
	assert.NotNil(t, ur)
	assert.Equal(t, http.StatusOK, ur.Status)
}
