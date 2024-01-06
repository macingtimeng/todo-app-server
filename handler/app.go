package handler

import (
	"fmt"

	"todo-app/handler/users_handler"
	"todo-app/infra/config"
	"todo-app/infra/db"
	"todo-app/repo/users_repo/users_pg"
	"todo-app/service/auth_service"
	"todo-app/service/users_service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartApp() {

	config.LoadEnv()
	db := db.DbConn()

	userRepo := users_pg.NewUsersRepo(db)
	userService := users_service.NewUserService(userRepo)
	userHandler := users_handler.NewUserHandler(userService)

	authService := auth_service.NewAuthService(userRepo)

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

	// users
	app.Post("/users/register", userHandler.Register)
	app.Post("/users/login", userHandler.Login)
	app.Get("/users/profile", authService.Authentication(), userHandler.Profile)
	app.Patch("/users/modify", authService.Authentication(), userHandler.Modify)

	app.Listen(":" + config.AppConfig().Port)
}
