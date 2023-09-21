package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id" gorm:"primaryKey"`
	Email    string    `json:"email" gorm:"unique"`
	Password []byte    `json:"-"`
	Alias    string    `json:"alias" gorm:"unique"`
	Name     string    `json:"name"`
}
