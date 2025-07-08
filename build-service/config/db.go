package config

import (
	"fmt"
	"log"
	"os"

	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/class"
	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/constants"
	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/skill"
	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/tag"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	// Importing for side effects - Dont Remove
	// This IS being used!
	_ "github.com/lib/pq"
)

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

	// pass the db connection string to connect to our database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	fmt.Printf("\nConnected to the database successfully.\n\n")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	SeedDefaults(db)

	return db
}

func runMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %v", err)
	}

	fmt.Printf("Successfully ran all migrations.\n\n")
	return nil
}

func SeedDefaults(db *sqlx.DB) {

	fmt.Printf("Successfully created all default members.\n\n")

	// --- default classes ---
	classRepo := class.NewRepository(db)
	classService := class.NewService(classRepo)

	err := classService.CreateDefaultClassesAndAscendancies(constants.DefaultClasses, constants.DefaultAscendancies)

	if err != nil {
		log.Fatal("Error when attempting to create default classes and ascendancies:", err)
	}

	fmt.Printf("Successfully created all default classes and ascendancies.\n\n")

	// --- default skills ---
	skillRepo := skill.NewRepository(db)
	skillService := skill.NewService(skillRepo)

	err = skillService.BatchCreateSkills(constants.ActiveSkills)
	if err != nil {
		log.Fatal("Error when attempting to create default active skills:", err)
	}

	err = skillService.BatchCreateSkills(constants.SupportSkills)
	if err != nil {
		log.Fatal("Error when attempting to create default support skills:", err)
	}

	// --- default tags ---
	tagsRepo := tag.NewRepository(db)
	tagsService := tag.NewService(tagsRepo)
	err = tagsService.CreateDefaultTags(constants.DefaultTags)

	// --- default items ---

	fmt.Printf("Successfully created all default active and support skills.\n\n")

	// itemRepo := item.NewItemRepository(db)
	// itemService := item.NewItemService(itemRepo, skillService)
	// itemService.CrawlingAndAddUniqueItemsService()
	// itemService.CrawlingAndAddBaseItemsService()
	// itemService.CrawlingAndAddItemModsService()
}
