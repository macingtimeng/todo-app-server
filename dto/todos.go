package dto

import (
	"time"
	"todo-app/entity"
)

type TodoResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type AddTodo struct {
	Todos string `json:"todos" valid:"required~ Todos can't be empty"`
}

type ModifyTodo struct {
	Todos  string `json:"todos" valid:"required~ Todos can't be empty"`
	Status bool   `json:"status"`
}

func (m *ModifyTodo) ModifyTodoToEntity() *entity.Todo {
	return &entity.Todo{
		Todos:  m.Todos,
		Status: m.Status,
	}
}

type Todo struct {
	Id        uint      `json:"id"`
	Todos     string    `json:"todos"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
