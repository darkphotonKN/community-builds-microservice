package config

import (
	"fmt"

	"log"
	"os"

	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/class"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/constants"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/member"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/skill"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/tag"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	// Importing for side effects - Dont Remove
	// This IS being used!
	_ "github.com/lib/pq"
)

// NOTE: for global db access, do not remove
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

	// seed default
	SeedDefaults(db)

	// set global instance for the database
	DB = db
	return DB
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
	// --- default members ---
	memberRepo := member.NewMemberRepository(db)
	memberService := member.NewMemberService(memberRepo)

	err := memberService.CreateDefaultMembersService(constants.DefaultMembers)

	if err != nil {
		log.Fatal("Error when attempting to create default members:", err)
	}

	fmt.Printf("Successfully created all default members.\n\n")

	// --- default classes ---
	classRepo := class.NewClassRepository(db)
	classService := class.NewClassService(classRepo)

	err = classService.CreateDefaultClassesAndAscendanciesService(constants.DefaultClasses, constants.DefaultAscendancies)

	if err != nil {
		log.Fatal("Error when attempting to create default classes and ascendancies:", err)
	}

	fmt.Printf("Successfully created all default classes and ascendancies.\n\n")

	// --- default skills ---
	skillRepo := skill.NewSkillRepository(db)
	skillService := skill.NewSkillService(skillRepo)

	err = skillService.BatchCreateSkillsService(constants.ActiveSkills)
	if err != nil {
		log.Fatal("Error when attempting to create default active skills:", err)
	}

	err = skillService.BatchCreateSkillsService(constants.SupportSkills)
	if err != nil {
		log.Fatal("Error when attempting to create default support skills:", err)
	}

	// --- default tags ---
	tagsRepo := tag.NewTagRepository(db)
	tagsService := tag.NewTagService(tagsRepo)
	err = tagsService.CreateDefaultTags(constants.DefaultTags)

	// --- default items ---

	fmt.Printf("Successfully created all default active and support skills.\n\n")

	// itemRepo := item.NewItemRepository(db)
	// itemService := item.NewItemService(itemRepo, skillService)
	// itemService.CrawlingAndAddUniqueItemsService()
	// itemService.CrawlingAndAddBaseItemsService()
	// itemService.CrawlingAndAddItemModsService()
}
