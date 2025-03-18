package database

import (
	"fmt"
	"log"

	"anchor-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
func MigrateAndSeedDatabase() (*gorm.DB, error) {
	// Initialize and migrate the database
	dsn := "host=anchor-backend_dev-db_1 user=docker password=docker dbname=docker port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
		return nil, err
	}

	// Drop existing tables
	err = db.Migrator().DropTable(&models.User{}, &models.Rack{})
	if err != nil {
		log.Fatal("Failed to drop existing tables: ", err)
		return nil, err
	}
	fmt.Println("Existing tables dropped successfully!")

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Rack{})
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
		return nil, err
	}
	fmt.Println("Database migrated successfully!")

	// Seed the database
	seedDatabase(db)

	return db, nil
}

// Seed function
func seedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		// Seed with initial users
		users := []models.User{
			{Name: "John Doe", Email: "john@example.com", Phone: "9059059050"},
			{Name: "Jane Smith", Email: "jane@example.com", Phone: "9059059051"},
			{Name: "test", Email: "evan@gmail.com", Phone: "test"},

		}

		racks := []models.Rack{
			{CurrUserID: 0},
			{CurrUserID: 2},
		}

		for _, user := range users {
			db.Create(&user)
		}

		for _, rack := range racks {
			db.Create(&rack)
		}

		fmt.Println("Database seeded successfully!")
	} else {
		fmt.Println("Database already seeded!")
	}
}
