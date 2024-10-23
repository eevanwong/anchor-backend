package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
// gorm.Model definition -> built in
type Model struct {
  ID        uint           `gorm:"primaryKey"` <- auto increments from 1
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}
*/
// Define your models
type User struct {
	gorm.Model
	Name  string
	Email string
}

type Rack struct {
	gorm.Model
	CurrUserID uint
}

// https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
func main() {
	// Initialize and migrate the database
	dsn := "host=localhost user=docker password=docker dbname=docker port=5434 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&User{}, &Rack{})
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}
	fmt.Println("Database migrated successfully!")

	// Seed the database
	seedDatabase(db)
}

// Seed function
func seedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		// Seed with initial users
		users := []User{
			{Name: "John Doe", Email: "john@example.com"},
			{Name: "Jane Smith", Email: "jane@example.com"},
		}

		racks := []Rack{
			{CurrUserID: 1},
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

