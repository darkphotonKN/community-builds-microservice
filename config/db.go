package config

import (
	"fmt"

	"log"
	"os"

	"github.com/darkphotonKN/community-builds/internal/class"
	"github.com/jmoiron/sqlx"

	// Importing for side effects - Dont Remove
	// This IS being used!
	_ "github.com/lib/pq"
)

// NOTE: for global db access, do not remove or move inside a function
var DB *sqlx.DB

/**
* Sets up the Database connection and provides its access as a singleton to
* the entire application.
**/
func InitDB() *sqlx.DB {
	// construct the db connection string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("constructed dsn", dsn)

	// pass the db connection string to connect to our database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	fmt.Println("Connected to the database successfully.")

	// seed default
	SeedDefaults(db)

	// set global instance for the database
	DB = db
	return DB
}

func SeedDefaults(db *sqlx.DB) {
	// default classes
	classes := []class.CreateClass{
		class.CreateClass{Name: "Warrior", Description: "Brutal monster wielding melee weapns.", ImageURL: "Placeholder."},
		class.CreateClass{Name: "Sorceror", Description: "Brutal monster wielding melee weapns.", ImageURL: "Placeholder."},
		class.CreateClass{Name: "Witch", Description: "Brutal monster wielding melee weapns.", ImageURL: "Placeholder."},
		class.CreateClass{Name: "Monk", Description: "Brutal monster wielding melee weapns.", ImageURL: "Placeholder."},
		class.CreateClass{Name: "Ranger", Description: "Brutal monster wielding melee weapns.", ImageURL: "Placeholder."},
		class.CreateClass{Name: "Mercenary", Description: "Brutal monster wielding melee weapns.", ImageURL: "Placeholder."},
	}
	classRepo := class.NewClassRepository(db)
	classService := class.NewClassService(classRepo)

	err := classService.CreateClassesAndAscendanciesService(classes)

	if err != nil {
		log.Fatalf("Error when attempting to create default classes and ascendancies:", err)
	}
}
