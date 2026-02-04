package config

import (
	"Lpanel/backend/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/lpanel.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(&models.Site{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
