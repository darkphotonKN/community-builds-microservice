package testsuite

import (
	"fmt"
	"log"

	"github.com/darkphotonKN/community-builds/config"
	"github.com/google/uuid"
)

/**
* Functions that create temporary fake data into the test database for unit and integration testing.
**/

func (t *TestSuite) seedTestData() {
	// seed defaults
	config.SeedDefaults(t.DB)

	//  test member and fetch its ID
	_, err := t.DB.Exec(`
		INSERT INTO members (name, email, password)
		VALUES ('Test Member', 'test.build@example.com', 'password')
		ON CONFLICT DO NOTHING
	`) // ON CONFLICT DO NOTHING - ignore insert if already created

	if err != nil {
		log.Fatalf("Failed to seed members: %v", err)
	}

	var memberID uuid.UUID

	err = t.DB.Get(&memberID, `
		SELECT id FROM members
		WHERE email = 'test.build@example.com'
	`)

	if err != nil {
		log.Fatalf("Failed to get member: %v", err)
	}

	fmt.Printf("Seeded member with ID: %s\n", memberID)

	// get primary skill for creating build
	var mainSkillId uuid.UUID
	err = t.DB.Get(&mainSkillId, `
		SELECT id FROM skills
		WHERE name = 'Earthquake'
	`)

	if err != nil {
		log.Fatalf("Failed to get skill with name '%s', error: %v", "Earthquake", err)
	}

	// get test class for creating build
	var classId uuid.UUID
	err = t.DB.Get(&classId, `
		SELECT 
			id
		FROM 
			classes
		WHERE name = 'Warrior'
	`)

	// insert test build using the memberID and one of the skills
	fmt.Printf("Seeding build with memberId: %s, mainSkillID: %s, classID: %s\n", memberID, mainSkillId, classId)

	testBuildName := "Earthquake Test Build"

	_, err = t.DB.Exec(`
		INSERT INTO builds (id, member_id, class_id, title, description, main_skill_id)
		VALUES (uuid_generate_v4(), $1, $3, $4, 'Description of Earthquake Test Build.',
		$2)
	`, memberID, classId, mainSkillId, testBuildName)

	if err != nil {
		log.Fatalf("Failed to seed build: %v", err)
	}

	// fetch and log build ID for verification

	var buildID uuid.UUID
	err = t.DB.Get(&buildID, `
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
		"memberId": memberID,
		"buildId":  buildID,
	})
}
