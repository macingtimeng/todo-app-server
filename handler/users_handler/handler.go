package users_handler

import (
	"todo-app/dto"
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/pkg/helper"
	"todo-app/service/users_service"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	us users_service.UserService
}

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Profile(c *fiber.Ctx) error
	Modify(c *fiber.Ctx) error
}

func NewUserHandler(userService users_service.UserService) UserHandler {
	return &userHandler{us: userService}
}

// Login implements UserHandler.
func (uh *userHandler) Login(c *fiber.Ctx) error {
	payload := &dto.Login{}

	if err := c.BodyParser(payload); err != nil {
		invalidJSON := errs.NewUnprocessableEntityError("invalid JSON body request")
		return c.Status(invalidJSON.Status()).JSON(invalidJSON)
	}

	err := helper.ValidateStruct(payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	ur, err := uh.us.Login(payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(ur.Status).JSON(ur)
}

// Modify implements UserHandler.
func (uh *userHandler) Modify(c *fiber.Ctx) error {
	payload := &dto.Modify{}
	user := c.Locals("user").(entity.User)

	if err := c.BodyParser(payload); err != nil {
		invalidJSON := errs.NewUnprocessableEntityError("invalid JSON body request")
		return c.Status(invalidJSON.Status()).JSON(invalidJSON)
	}

	err := helper.ValidateStruct(payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	ur, err := uh.us.Modify(user.ID, payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(ur.Status).JSON(ur)
}

// Profile implements UserHandler.
func (uh *userHandler) Profile(c *fiber.Ctx) error {
	user := c.Locals("user").(entity.User)

	ur, err := uh.us.Profile(user.ID)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(ur.Status).JSON(ur)
}

// Register implements UserHandler.
func (uh *userHandler) Register(c *fiber.Ctx) error {
	payload := &dto.Register{}

	if err := c.BodyParser(payload); err != nil {
		invalidJSON := errs.NewUnprocessableEntityError("invalid JSON body request")
		return c.Status(invalidJSON.Status()).JSON(invalidJSON)
	}

	err := helper.ValidateStruct(payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	ur, err := uh.us.Register(payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(ur.Status).JSON(ur)
}
