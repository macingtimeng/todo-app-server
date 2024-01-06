package auth_service

import fiber "github.com/gofiber/fiber/v2"

type authMock struct {
}

var (
	Authentication func() fiber.Handler
)

func NewAuthMock() AuthService {
	return &authMock{}
}

// Authentication implements AuthService.
func (a *authMock) Authentication() fiber.Handler {
	return Authentication()
}
