package dto

import "todo-app/entity"

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Token struct {
	TokenString string `json:"token"`
}

type Register struct {
	Name     string `json:"name" valid:"required~ Name can't be empty"`
	Email    string `json:"email" valid:"required~ Email can't be empty, email"`
	Password string `json:"password" valid:"required~ Password can't be empty"`
}

func (r *Register) RegisterToEntity() *entity.User {
	return &entity.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

type Login struct {
	Email    string `json:"email" valid:"required~ Email can't be empty, email"`
	Password string `json:"password" valid:"required~ Password can't be empty"`
}

type Modify struct {
	Name  string `json:"name" valid:"required~ Name can't be empty"`
	Email string `json:"email" valid:"required~ Email can't be empty, email"`
}

func (m *Modify) ModifyToEntity() *entity.User {
	return &entity.User{
		Name:  m.Name,
		Email: m.Email,
	}
}

type Profile struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
