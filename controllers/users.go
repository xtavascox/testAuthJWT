package controllers

import (
	"github.com/authentication-app/repository"
	"github.com/gofiber/fiber/v2"
)

func Users(app *fiber.App) {
	user := app.Group("/api/v1/user")

	user.Post("/register", repository.RegisterUser)
	user.Post("/login", repository.LoginUser)
	user.Post("/logout", repository.Logout)
	user.Get("/:id", repository.User)
}
