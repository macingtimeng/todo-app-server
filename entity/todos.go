package entity

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Todos  string
	Status bool
	UserID uint
}
