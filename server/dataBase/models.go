package dataBase

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UserSession struct {
	gorm.Model
	SessionKey string `json:"session_key"`
	UserId     uint   `json:"user_id"`
	User       User   `gorm:"references:ID"`
}
