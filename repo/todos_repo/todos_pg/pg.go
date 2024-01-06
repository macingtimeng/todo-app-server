package todos_pg

import (
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/repo/todos_repo"

	"gorm.io/gorm"
)

type todoPg struct {
	db *gorm.DB
}

func NewTodoRepo(db *gorm.DB) todos_repo.TodoRepo {
	return &todoPg{db: db}
}

// Add implements todos_repo.TodoRepo.
func (pg *todoPg) Add(todo *entity.Todo) errs.Error {
	tx := pg.db.Begin()

	if err := tx.Create(todo).Error; err != nil {
		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}

// Delete implements todos_repo.TodoRepo.
func (pg *todoPg) Delete(todoId uint) errs.Error {

	tx := pg.db.Begin()

	if err := tx.Delete(&entity.Todo{}, "id = ?", todoId).Error; err != nil {
		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}

// Detail implements todos_repo.TodoRepo.
func (pg *todoPg) Detail(todoId uint) (*entity.Todo, errs.Error) {

	todo := entity.Todo{}

	if err := pg.db.First(&todo, "id = ?", todoId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFoundError("todo not found")
		}
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &todo, nil
}

// Fetch implements todos_repo.TodoRepo.
func (pg *todoPg) Fetch(userId uint) ([]*entity.Todo, errs.Error) {

	todos := []*entity.Todo{}

	if err := pg.db.Find(&todos, "user_id = ?", userId).Error; err != nil {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return todos, nil
}

// Modify implements todos_repo.TodoRepo.
func (pg *todoPg) Modify(todoId uint, todo *entity.Todo) errs.Error {

	tx := pg.db.Begin()

	if err := tx.Model(&entity.Todo{}).Where("id = ?", todoId).Updates(todo).Error; err != nil {
		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}
