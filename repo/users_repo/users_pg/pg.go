package users_pg

import (
	"strings"
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/repo/users_repo"

	"gorm.io/gorm"
)

type usersPg struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) users_repo.UsersRepo {
	return &usersPg{db: db}
}

// Add implements users_repo.UsersRepo.
func (pg *usersPg) Add(user *entity.User) errs.Error {

	tx := pg.db.Begin()

	if err := tx.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			tx.Rollback()
			return errs.NewConflictError("email has been used")
		}

		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}

// FetchByEmail implements users_repo.UsersRepo.
func (pg *usersPg) FetchByEmail(email string) (*entity.User, errs.Error) {

	users := entity.User{}

	if err := pg.db.First(&users, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFoundError("user not found")
		}
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &users, nil
}

// FetchById implements users_repo.UsersRepo.
func (pg *usersPg) FetchById(userId uint) (*entity.User, errs.Error) {

	users := entity.User{}

	if err := pg.db.First(&users, "id = ?", userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFoundError("user not found")
		}
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &users, nil
}

// Modify implements users_repo.UsersRepo.
func (pg *usersPg) Modify(userId uint, user *entity.User) errs.Error {

	tx := pg.db.Begin()

	if err := tx.Model(&entity.User{}).Where("id = ?", userId).Updates(user).Error; err != nil {
		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}
