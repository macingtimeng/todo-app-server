package todos_handler

import (
	"strconv"
	"todo-app/dto"
	"todo-app/entity"
	"todo-app/pkg/errs"
	"todo-app/pkg/helper"
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
// Add godoc
// @Summary Add todo
// @Description Add todo request
// @Tags Todos
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param dto.AddTodo body dto.AddTodo true "body request for add todo"
// @Success 201 {object} dto.TodoResponse
// @Router /todos [post]
func (th *todoHandler) Add(c *fiber.Ctx) error {
	payload := &dto.AddTodo{}
	user := c.Locals("user").(entity.User)

	if err := c.BodyParser(payload); err != nil {
		invalidJson := errs.NewUnprocessableEntityError("invalid JSON body request")
		return c.Status(invalidJson.Status()).JSON(invalidJson)
	}

	if err := helper.ValidateStruct(payload); err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	tr, err := th.ts.Add(user.ID, payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}

// Delete implements TodoHandler.
// Delete godoc
// @Summary Delete todo
// @Description Delete todo request
// @Tags Todos
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param todoId path int true "todo id"
// @Success 200 {object} dto.TodoResponse
// @Router /todos/{todoId} [delete]
func (th *todoHandler) Delete(c *fiber.Ctx) error {

	todoId, _ := strconv.Atoi(c.Params("todoId"))

	tr, err := th.ts.Delete(uint(todoId))

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}

// Detail implements TodoHandler.
// Detail godoc
// @Summary Detail todo
// @Description Detail todo request
// @Tags Todos
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param todoId path int true "todo id"
// @Success 200 {object} dto.TodoResponse
// @Router /todos/{todoId} [get]
func (th *todoHandler) Detail(c *fiber.Ctx) error {

	todoId, _ := strconv.Atoi(c.Params("todoId"))

	tr, err := th.ts.Detail(uint(todoId))

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}

// Fetch implements TodoHandler.
// Fetch godoc
// @Summary Get all todos
// @Description Get all todos request
// @Tags Todos
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} dto.TodoResponse
// @Router /todos/ [get]
func (th *todoHandler) Fetch(c *fiber.Ctx) error {

	user := c.Locals("user").(entity.User)

	tr, err := th.ts.Fetch(user.ID)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}

// Modify implements TodoHandler.
// Modify godoc
// @Summary Modify todo
// @Description Modify todo request
// @Tags Todos
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param todoId path int true "todo id"
// @Param dto.ModifyTodo body dto.ModifyTodo true "body request for modify todo"
// @Success 200 {object} dto.TodoResponse
// @Router /todos/{todoId} [patch]
func (th *todoHandler) Modify(c *fiber.Ctx) error {
	payload := &dto.ModifyTodo{}
	todoId, _ := strconv.Atoi(c.Params("todoId"))

	if err := c.BodyParser(payload); err != nil {
		invalidJson := errs.NewUnprocessableEntityError("invalid JSON body request")
		return c.Status(invalidJson.Status()).JSON(invalidJson)
	}

	if err := helper.ValidateStruct(payload); err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	tr, err := th.ts.Modify(uint(todoId), payload)

	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	return c.Status(tr.Status).JSON(tr)
}
