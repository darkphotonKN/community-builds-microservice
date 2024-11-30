package build_test

import (
	"fmt"
	"testing"

	"github.com/darkphotonKN/community-builds/internal/build"
	"github.com/darkphotonKN/community-builds/internal/skill"
	test "github.com/darkphotonKN/community-builds/internal/test/setup"
	"github.com/google/uuid"
)

func TestAddSkillLinksToBuildService_Success(t *testing.T) {
	// -- setup and seed test database --
	testSuite := test.NewTestSuite()

	testDB := testSuite.SetupTestDB()

	// test build name
	testBuildName := "Earthquake Test Build"

	// get a member's id
	metadata, metadataOk := testSuite.Metadata["TestAddSkillLinksToBuildService_Success"].(map[string]interface{})

	mainSkillId, skillIdOk := metadata["mainSkillId"].(uuid.UUID)
	memberId, memberIdOk := metadata["memberId"].(uuid.UUID)

	if !metadataOk || !skillIdOk || !memberIdOk {
		panic("mainSkillId assertion failed.")
	}

	fmt.Printf("\n@TEST: assert mainSkillId: %s and memberId: %s\n\n", memberId, mainSkillId)

	// get the build id responsible for testing
	var buildId uuid.UUID

	testDB.Get(&buildId, `
		SELECT id FROM builds
		WHERE title = $1 
		`, testBuildName)

	fmt.Println("@TEST: retrieved buildId:", buildId)

	// -- test service methods --

	// --- DI setup ---
	skillRepo := skill.NewSkillRepository(testDB)
	skillService := skill.NewSkillService(skillRepo)

	buildsRepo := build.NewBuildRepository(testDB)
	buildsService := build.NewBuildService(buildsRepo, skillService)

	fmt.Println(buildsService)

	// --- tests ---

	// mock payload
	// payload := build.AddSkillsToBuildRequest{
	// MainSkillLinks: build.SkillLinks{
	// 	SkillLinkName: "Earthquake",
	// 	Skill:         "",
	// 	Links,
	// },
	// }

	// buildsService.AddSkillLinksToBuildService(memberId, buildId, payload)
}
