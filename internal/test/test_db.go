package testsuite

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	// Importing for side effects - Dont Remove
	// This IS being used!
	_ "github.com/lib/pq"
)

func (t *TestSuite) ConnectTestDB() *sqlx.DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_NAME"),
	)

	fmt.Println("constructed dsn", dsn)

	// pass the db connection string to connect to our database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	fmt.Println("Connected to the database successfully.")
	return db
}

/**
* Resets database.
**/
func (t *TestSuite) ResetTestDB(db *sqlx.DB) {
	// Drop all tables and reset the schema
	_, err := db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
	if err != nil {
		log.Fatalf("Failed to reset test database: %v", err)
	}
}

/**
* Creates test tables for Test DB.
**/
func (t *TestSuite) applyTestMigrations() {
	migrations := []string{
		// uuid-ossp extension for uuid v4 generation
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,

		// members table
		`CREATE TABLE IF NOT EXISTS members (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			status SMALLINT NOT NULL DEFAULT 1 CHECK (status IN (1, 2)),
			average_rating DECIMAL(2, 1) DEFAULT 0 CHECK (average_rating >= 0 AND average_rating <= 5)
		);`,

		// skills table
		`
		CREATE TABLE IF NOT EXISTS skills (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name TEXT NOT NULL UNIQUE,  
			type TEXT NOT NULL CHECK (type IN ('active', 'support')), 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,

		// builds table
		`
		CREATE TABLE IF NOT EXISTS builds (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			member_id UUID NOT NULL, -- FK to members table
			title TEXT NOT NULL, 
			description TEXT NOT NULL,
			main_skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE RESTRICT, 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,

		// build Skill Links table
		`
		CREATE TABLE IF NOT EXISTS build_skill_links (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,  
			name TEXT NOT NULL,
			is_main BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,

		// build Skill Link Skills table
		`
		CREATE TABLE IF NOT EXISTS build_skill_link_skills (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			build_skill_link_id UUID NOT NULL REFERENCES build_skill_links(id) ON DELETE CASCADE,
			skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE
		);
		`,
	}

	for _, migration := range migrations {
		_, err := t.DB.Exec(migration)
		if err != nil {
			log.Fatalf("Failed to apply test migration: %v", err)
		}
	}

	fmt.Println("Test database migrations applied successfully.")
}
