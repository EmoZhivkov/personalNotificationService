package repositories

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getTestInMemoryDBClient() *gorm.DB {
	dbClient, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := dbClient.AutoMigrate(tablesToCreate...); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return dbClient
}
