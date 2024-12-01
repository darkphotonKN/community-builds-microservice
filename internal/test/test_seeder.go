package testsuite

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

/**
* Functions that create temporary fake data into the test database for unit and integration testing.
**/

func (t *TestSuite) seedTestData(db *sqlx.DB) {
	//  test member and fetch its ID
	_, err := db.Exec(`
		INSERT INTO members (name, email, password)
		VALUES ('Test Member', 'test.build@example.com', 'password')
		ON CONFLICT DO NOTHING 
	`) // ON CONFLICT DO NOTHING - ignore insert if already created

	if err != nil {
		log.Fatalf("Failed to seed members: %v", err)
	}

	var memberID uuid.UUID

	err = db.Get(&memberID, `
		SELECT id FROM members
		WHERE email = 'test.build@example.com'
	`)

	if err != nil {
		log.Fatalf("Failed to get member: %v", err)
	}

	fmt.Printf("Seeded member with ID: %s\n", memberID)

	earthquakeId, _ := uuid.Parse("8ce22ebc-e367-4b98-a48b-49f984d8fb8d")
	lightningStrikeId, _ := uuid.Parse("2a3810ee-cb55-4656-83be-2c3d14683495")
	increasedAOEId, _ := uuid.Parse("23a87914-ffc0-4f73-b113-41c361e740ee")
	multistrikeId, _ := uuid.Parse("5f2eb204-40af-4bce-9ffd-296fcebac9b6")
	concentratedEffectId, _ := uuid.Parse("92618534-d832-4afd-a01b-d453ffa5658d")

	// NOTE: transferring meta data for tests
	// ["name of test"] -> map of data needed

	// create test skills
	skills := []struct {
		ID   uuid.UUID
		Name string
		Type string
	}{
		{lightningStrikeId, "Lightning Strike", "active"},
		{earthquakeId, "Earthquake", "active"},
		{increasedAOEId, "Increased Area of Effect", "support"},
		{multistrikeId, "Multistrike", "support"},
		{concentratedEffectId, "Concentrated Effect", "support"},
	}

	for _, skill := range skills {
		_, err = db.Exec(`
			INSERT INTO skills (id, name, type)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`, skill.ID, skill.Name, skill.Type) // ON CONFLICT DO NOTHING - ignore insert if already created

		if err != nil {
			log.Fatalf("Failed to seed skill '%s': %v", skill.Name, err)
		}
	}

	// get primary skill for creating build
	var mainSkillId uuid.UUID
	err = db.Get(&mainSkillId, `
		SELECT id FROM skills
		WHERE name = 'Earthquake'
	`)

	if err != nil {
		log.Fatalf("Failed to get skill with name '%s', error: %v", "Earthquake", err)
	}

	// insert test build using the memberID and one of the skills
	fmt.Printf("Seeding build with memberId: %s, mainSkillID: %s\n", memberID, mainSkillId)

	testBuildName := "Earthquake Test Build"

	_, err = db.Exec(`
		INSERT INTO builds (id, member_id, title, description, main_skill_id)
		VALUES (uuid_generate_v4(), $1, $3, 'Description of Earthquake Test Build.', 
		$2)
	`, memberID, mainSkillId, testBuildName) // ON CONFLICT DO NOTHING - ignore insert if already created
	if err != nil {
		log.Fatalf("Failed to seed build: %v", err)
	}

	// fetch and log build ID for verification

	var buildID uuid.UUID
	err = db.Get(&buildID, `
		SELECT id FROM builds
		WHERE title = $1
	`, testBuildName)

	if err != nil {
		log.Fatalf("Failed to fetch build ID: %v", err)
	}

	fmt.Printf("Seeded build with ID: %s\n", buildID)

	// --- META DATA ---

	// -- TestAddSkillLinksToBuildService_Success METADATA --
	t.setMetaData("TestAddSkillLinksToBuildService_Success", MetaData{
		"memberId":    memberID,
		"mainSkillId": earthquakeId,
		"buildId":     buildID,
	})

}
