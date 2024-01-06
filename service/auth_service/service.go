package auth_service

import (
	"strconv"
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/repo/todos_repo"
	"todo-app/repo/users_repo"

	"github.com/gofiber/fiber/v2"
)

type authService struct {
	ur users_repo.UsersRepo
	tr todos_repo.TodoRepo
}

type AuthService interface {
	Authentication() fiber.Handler
	Authorization() fiber.Handler
}

func NewAuthService(userRepo users_repo.UsersRepo, todoRepo todos_repo.TodoRepo) AuthService {
	return &authService{ur: userRepo, tr: todoRepo}
}

// Authentication implements AuthService.
func (as *authService) Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {

		bearerToken := c.Get("Authorization")
		user := entity.User{}

		err := user.ValidateToken(bearerToken)

		if err != nil {
			return c.Status(err.Status()).JSON(err)
		}

		_, err = as.ur.FetchByEmail(user.Email)

		if err != nil {
			errUnauthenticated := errs.NewUnauthenticatedError("invalid user")
			return c.Status(errUnauthenticated.Status()).JSON(errUnauthenticated)
		}

		c.Locals("user", user)

		return c.Next()
	}
}

// Authorization implements AuthService.
func (as *authService) Authorization() fiber.Handler {
	return func(c *fiber.Ctx) error {

		user := c.Locals("user").(entity.User)
		todoId, _ := strconv.Atoi(c.Params("todoId"))

		t, err := as.tr.Detail(uint(todoId))

		if err != nil {
			return c.Status(err.Status()).JSON(err)
		}

		if t.UserID != user.ID {
			errUnauthorizedError := errs.NewUnathorizedError("you're not authorized to access this todo")
			return c.Status(errUnauthorizedError.Status()).JSON(errUnauthorizedError)
		}

		return c.Next()
	}
}
