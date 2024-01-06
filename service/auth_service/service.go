package auth_service

import (
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/repo/users_repo"

	"github.com/gofiber/fiber/v2"
)

type authService struct {
	ur users_repo.UsersRepo
}

type AuthService interface {
	Authentication() fiber.Handler
}

func NewAuthService(userRepo users_repo.UsersRepo) AuthService {
	return &authService{ur: userRepo}
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
