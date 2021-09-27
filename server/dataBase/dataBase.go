package dataBase

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase() {
	connectionString := "host=localhost port=5432 user=postgres dbname=awesomeDB password=1 sslmode=disable"

	database, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к БД")
	}

	database.AutoMigrate(&User{}, &UserSession{})

	DB = database
}
