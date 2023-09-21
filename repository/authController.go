package repository

import (
	"github.com/authentication-app/database"
	"github.com/authentication-app/models"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const SecretKey = "secret"

func Hello(ctx *fiber.Ctx) error {

	return ctx.SendString("Hello, World from users 👋!")
}
func RegisterUser(ctx *fiber.Ctx) error {
	var data map[string]string
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Id:       uuid.New(),
		Email:    data["email"],
		Password: password,
		Alias:    data["user_name"],
		Name:     data["name"],
	}
	database.DB.Create(&user)

	return ctx.JSON(user)
}

func LoginUser(ctx *fiber.Ctx) error {
	var data map[string]string
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}
	log.Println("entro a login")
	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == uuid.Nil {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if error := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); error != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}
	// expiration time of 24 hours to use in claims
	expireAt := time.Unix(time.Now().Add(time.Minute*5).Unix(), 0)

	//claims := jwt.MapClaims{
	//	"email": user.Email,
	//	"alias": user.Alias,
	//	"iss":   user.Id.String(),
	//	"exp":   jwt.At(expireAt),
	//}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Id.String(),
		ExpiresAt: jwt.At(expireAt),
	})
	//newWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(fiber.Map{
			"message": "Could not login",
		})
	}

	cookie := fiber.Cookie{Name: "jwt", Value: token, Expires: expireAt, HTTPOnly: true}

	ctx.Cookie(&cookie)

	return ctx.JSON(fiber.Map{
		"message": "Success",
		"token":   token,
	})
}

func User(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	log.Println(id)
	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user); err.Error != nil {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Success",
		"user":    user,
	})
}

func Logout(ctx *fiber.Ctx) error {
	cookie := fiber.Cookie{Name: "jwt", Value: "", Expires: time.Now().Add(-time.Hour), HTTPOnly: true}

	ctx.Cookie(&cookie)

	return ctx.JSON(fiber.Map{
		"message": "Success",
	})
}

func ValidateJwt(c *fiber.Ctx) error {
	if c.Path() == "/api/v1/user/login" {
		return c.Next()
	}
	log.Println("entro a validar")
	log.Println(c.Path())
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	return c.Next()
}
