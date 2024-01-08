package handler

import (
	"fmt"

	"todo-app/handler/todos_handler"
	"todo-app/handler/users_handler"
	"todo-app/infra/config"
	"todo-app/infra/db"
	"todo-app/repo/todos_repo/todos_pg"
	"todo-app/repo/users_repo/users_pg"
	"todo-app/service/auth_service"
	"todo-app/service/todos_service"
	"todo-app/service/users_service"

	"github.com/gofiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title TodoKu API V1
// @version 1.0
// @description 

// @contact.name Yusril Ilham Kholid
// @contact.email yusrililham62@gmail.com
// @contact.url https://yusrililhm.netlify.app

// @host localhost:8080
// @BasePath /api/v1/
func StartApp() {

	config.LoadEnv()
	db := db.DbConn()

	userRepo := users_pg.NewUsersRepo(db)
	userService := users_service.NewUserService(userRepo)
	userHandler := users_handler.NewUserHandler(userService)

	todoRepo := todos_pg.NewTodoRepo(db)
	todoService := todos_service.NewTodoService(todoRepo)
	todoHandler := todos_handler.NewTodoHandler(todoService)

	authService := auth_service.NewAuthService(userRepo, todoRepo)

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: fmt.Sprintf(
			"%s, %s, %s, %s, %s",
			fiber.MethodPost,
			fiber.MethodGet,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		),
		AllowHeaders: "Content-Type, Authorization",
	}))

	// swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// users
	app.Post("/api/v1/users/register", userHandler.Register)
	app.Post("/api/v1/users/login", userHandler.Login)
	app.Get("/api/v1/users/profile", authService.Authentication(), userHandler.Profile)
	app.Patch("/api/v1/users/modify", authService.Authentication(), userHandler.Modify)

	app.Post("/api/v1/todos", authService.Authentication(), todoHandler.Add)
	app.Get("/api/v1/todos", authService.Authentication(), todoHandler.Fetch)
	app.Delete("/api/v1/todos/:todoId", authService.Authentication(), authService.Authorization(), todoHandler.Delete)
	app.Get("/api/v1/todos/:todoId", authService.Authentication(), authService.Authorization(), todoHandler.Detail)
	app.Patch("/api/v1/todos/:todoId", authService.Authentication(), authService.Authorization(), todoHandler.Modify)

	app.Listen(":" + config.AppConfig().Port)
}
