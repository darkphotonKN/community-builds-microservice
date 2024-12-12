package config

import (
	"fmt"

	"log"
	"os"

	"github.com/darkphotonKN/community-builds/internal/class"
	"github.com/darkphotonKN/community-builds/internal/constants"
	"github.com/darkphotonKN/community-builds/internal/member"
	"github.com/darkphotonKN/community-builds/internal/skill"
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

	// seed default
	SeedDefaults(db)

	// set global instance for the database
	DB = db
	return DB
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

	fmt.Printf("Successfully created all default active and support skills.\n\n")

}
