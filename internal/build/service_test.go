package build_test

import (
	"fmt"
	"testing"

	"github.com/darkphotonKN/community-builds/internal/build"
	"github.com/darkphotonKN/community-builds/internal/skill"
	testsuite "github.com/darkphotonKN/community-builds/internal/test"
)

func TestAddSkillLinksToBuildService_Success(t *testing.T) {
	// --- setup and seed test database ---

	testSuite := testsuite.NewTestSuite()

	testDB := testSuite.SetupTestDB()

	// -- retrieve stored metadata --
	metadata, metadataOk := testSuite.Metadata["TestAddSkillLinksToBuildService_Success"].(map[string]interface{}) // global metadata
	mainSkillId, skillIdOk := metadata["mainSkillId"].(uuid.UUID)                                                  // get main skill id
	memberId, memberIdOk := metadata["memberId"].(uuid.UUID)                                                       // get a member's id
	buildId, buildIdOk := metadata["buildId"].(uuid.UUID)                                                          // get buildId

	if !metadataOk || !skillIdOk || !memberIdOk || !buildIdOk {
		t.Fatalf("Failed to retrieve metadata: metadataOk=%v, skillIdOk=%v, memberIdOk=%v", metadataOk, skillIdOk, memberIdOk)
	}

	fmt.Printf("\n@TEST: mainSkillId: %s\nmemberId: %s\nbuildId: %s\n\n", memberId, mainSkillId, buildId)

	// --- test service methods ---

	// -- DI setup --
	skillRepo := skill.NewSkillRepository(testDB)
	skillService := skill.NewSkillService(skillRepo)

	buildsRepo := build.NewBuildRepository(testDB)
	buildsService := build.NewBuildService(buildsRepo, skillService)

	fmt.Println(buildsService)

	// -- tests --

	// mock payload
	payload := build.AddSkillsToBuildRequest{
		MainSkillLinks: build.SkillLinks{
			SkillLinkName: "Earthquake",
			Skill:         mainSkillId,
		},
	}

	buildsService.AddSkillLinksToBuildService(memberId, buildId, payload)
}
