package db

import (
	"go-ecommerce/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDB() {
	dsn := "root:@tcp(localhost:3306)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to Database %v", err)
	}

	if DB == nil {
		log.Fatalf("Database connection is nil")
	}

	err = DB.Statement.AutoMigrate(&models.Product{})
	if err != nil {
		log.Printf("Failed to Auto-Migrate %v", err)
	}

}
