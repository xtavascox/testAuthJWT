package main

import (
	"github.com/authentication-app/controllers"
	"github.com/authentication-app/database"
	"github.com/authentication-app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	database.Connect()
	app := fiber.New()

	app.Use(repository.ValidateJwt)
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	controllers.Users(app)

	log.Fatal(app.Listen(":4000"))
}
