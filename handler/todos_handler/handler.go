package todos_handler

import (
	"strconv"
	"todo-app/dto"
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/service/todos_service"

	"github.com/gofiber/fiber/v2"
)

type todoHandler struct {
	ts todos_service.TodoService
}

type TodoHandler interface {
	Add(c *fiber.Ctx) error
	Detail(c *fiber.Ctx) error
	Fetch(c *fiber.Ctx) error
	Modify(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func NewTodoHandler(todoService todos_service.TodoService) TodoHandler {
	return &todoHandler{ts: todoService}
}

// Add implements TodoHandler.
func (th *todoHandler) Add(c *fiber.Ctx) error {
	payload := &dto.AddTodo{}
	user := c.Locals("user").(entity.User)

	if err := c.BodyParser(payload); err != nil {
		invalidJson := errs.NewUnprocessableEntityError("invalid JSON body request")
		return c.Status(invalidJson.Status()).JSON(invalidJson)
	}

	tr, err := th.ts.Add(user.ID, payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}

// Delete implements TodoHandler.
func (th *todoHandler) Delete(c *fiber.Ctx) error {

	todoId, _ := strconv.Atoi(c.Params("todoId"))

	tr, err := th.ts.Delete(uint(todoId))

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}

// Detail implements TodoHandler.
func (th *todoHandler) Detail(c *fiber.Ctx) error {

	todoId, _ := strconv.Atoi(c.Params("todoId"))

	tr, err := th.ts.Detail(uint(todoId))

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}

// Fetch implements TodoHandler.
func (th *todoHandler) Fetch(c *fiber.Ctx) error {

	user := c.Locals("user").(entity.User)

	tr, err := th.ts.Fetch(user.ID)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}

// Modify implements TodoHandler.
func (th *todoHandler) Modify(c *fiber.Ctx) error {
	payload := &dto.ModifyTodo{}
	todoId, _ := strconv.Atoi(c.Params("todoId"))

	if err := c.BodyParser(payload); err != nil {
		invalidJson := errs.NewUnprocessableEntityError("invalid JSON body request")
		return c.Status(invalidJson.Status()).JSON(invalidJson)
	}

	tr, err := th.ts.Modify(uint(todoId), payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}
