package database

import (
	"github.com/authentication-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect() {
	conn, err := gorm.Open(postgres.Open("host=localhost user=postgres password=admin dbname=test_auth port=5432 sslmode=disable TimeZone=America/Bogota"), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo conectar a la base de datos", err)
	}
	log.Println("conectado a la base de datos", conn)
	DB = conn
	conn.AutoMigrate(&models.User{})
}
