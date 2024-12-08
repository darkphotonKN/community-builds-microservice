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

/**
* Establishes the initial connection to the test database
**/
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

		// Members table
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

		// Skills table
		`CREATE TABLE IF NOT EXISTS skills (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name TEXT NOT NULL UNIQUE,
        type TEXT NOT NULL CHECK (type IN ('active', 'support')),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`,

		// Base Items table
		`CREATE TABLE IF NOT EXISTS base_items (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        category TEXT NOT NULL,
        class TEXT NOT NULL,
        type TEXT NOT NULL,
        name TEXT NOT NULL,
        equip_type TEXT NOT NULL,
        is_two_hands BOOLEAN NOT NULL,
        slot TEXT NOT NULL,
        image_url TEXT,
        required_level TEXT,
        required_strength TEXT,
        required_dexterity TEXT,
        required_intelligence TEXT,
        armour TEXT,
        energy_shield TEXT,
        evasion TEXT,
        ward TEXT,
        damage TEXT,
        aps TEXT,
        crit TEXT,
        dps TEXT,
        implicit TEXT[],
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`,

		// Items table
		`CREATE TABLE IF NOT EXISTS items (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        member_id UUID REFERENCES members(id) ON DELETE RESTRICT,
        base_item_id UUID REFERENCES base_items(id) ON DELETE RESTRICT,
        category TEXT NOT NULL,
        class TEXT NOT NULL,
        type TEXT NOT NULL,
        name TEXT NOT NULL,
        unique_item BOOLEAN NOT NULL,
        slot TEXT NOT NULL,
        description TEXT,
        image_url TEXT,
        required_level TEXT,
        required_strength TEXT,
        required_dexterity TEXT,
        required_intelligence TEXT,
        armour TEXT,
        block TEXT,
        energy_shield TEXT,
        evasion TEXT,
        ward TEXT,
        damage TEXT,
        aps TEXT,
        crit TEXT,
        pdps TEXT,
        edps TEXT,
        dps TEXT,
        life TEXT,
        mana TEXT,
        duration TEXT,
        usage TEXT,
        capacity TEXT,
        additional TEXT,
        stats TEXT[],
        implicit TEXT[],
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`,

		// Classes table
		`CREATE TABLE IF NOT EXISTS classes (
          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
          name TEXT NOT NULL UNIQUE,
          description TEXT NOT NULL,
          image_url TEXT,
          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      );`,

		// Ascendancies table
		`CREATE TABLE IF NOT EXISTS ascendancies (
          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
          class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
          name TEXT NOT NULL UNIQUE,
          description TEXT,
          image_url TEXT,
          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      );`,

		// Builds table
		`CREATE TABLE IF NOT EXISTS builds (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
			member_id UUID NOT NULL REFERENCES members (id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			main_skill_id UUID NOT NULL REFERENCES skills (id) ON DELETE RESTRICT,
			class_id UUID NOT NULL REFERENCES classes (id) ON DELETE RESTRICT,
			ascendancy_id UUID REFERENCES ascendancies (id) ON DELETE RESTRICT,
			avg_end_game_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
				avg_end_game_rating >= 0
				AND avg_end_game_rating <= 10
			),
			avg_fun_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
				avg_fun_rating >= 0
				AND avg_fun_rating <= 10
			),
			avg_creative_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
				avg_creative_rating >= 0
				AND avg_creative_rating <= 10
			),
			avg_speed_farm_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
				avg_speed_farm_rating >= 0
				AND avg_speed_farm_rating <= 10
			),
			avg_bossing_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
				avg_bossing_rating >= 0
				AND avg_bossing_rating <= 10
			),
			views INT DEFAULT 0 CHECK (views >= 0),
			status SMALLINT NOT NULL DEFAULT 0 CHECK (status IN (0, 1, 2)), -- 0: Edit, 1: Published, 2: Archived
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,

		// Ratings table
		`CREATE TABLE IF NOT EXISTS ratings (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
        member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
        category TEXT NOT NULL CHECK (
            category IN ('endgame', 'fun', 'creative', 'speedfarm', 'bossing')
        ),
        value SMALLINT CHECK (value BETWEEN 1 AND 10) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`,

		// Build items table
		`CREATE TABLE IF NOT EXISTS build_items (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
        item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
        slot TEXT NOT NULL
    );`,

		// Tags table
		`CREATE TABLE IF NOT EXISTS tags (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name TEXT NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`,

		// Build tags table
		`CREATE TABLE IF NOT EXISTS build_tags (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
        build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
        UNIQUE(build_id, tag_id)
    );`,

		// Build skill links table
		`CREATE TABLE IF NOT EXISTS build_skill_links (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        is_main BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`,

		// Build skill link skills table
		`CREATE TABLE IF NOT EXISTS build_skill_link_skills (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        build_skill_link_id UUID NOT NULL REFERENCES build_skill_links(id) ON DELETE CASCADE,
        skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE
    );`}

	for _, migration := range migrations {
		_, err := t.DB.Exec(migration)
		if err != nil {
			log.Fatalf("Failed to apply test migration: %v", err)
		}
	}

	fmt.Println("Test database migrations applied successfully.")
}
