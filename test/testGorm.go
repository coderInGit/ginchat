package main

import (
	"ginChat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:1234@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"),
		&gorm.Config{})
	if err != nil {
		panic("failed to connect databases")
	}
	db.AutoMigrate(&models.UserBasic{})

	user := &models.UserBasic{}
	user.Name = "申专"
	db.Create(user)

	db.First(&user, 1)
	db.First(&user, "code = ?", "D42")
}
