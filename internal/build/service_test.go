package build_test

import (
	"fmt"
	"testing"

	"github.com/darkphotonKN/community-builds/internal/build"
	"github.com/darkphotonKN/community-builds/internal/skill"
	testsuite "github.com/darkphotonKN/community-builds/internal/test"
	"github.com/google/uuid"
)

func TestAddSkillLinksToBuildService_Success(t *testing.T) {
	// --- setup and seed test database ---

	newTestSuite := testsuite.NewTestSuite()

	testDB := newTestSuite.SetupTestDB()

	// -- retrieve stored metaData --
	metaData := newTestSuite.GetMetaData("TestAddSkillLinksToBuildService_Success") // global metadata
	mainSkillId, skillIdOk := metaData["mainSkillId"].(uuid.UUID)                   // get main skill id
	memberId, memberIdOk := metaData["memberId"].(uuid.UUID)                        // get a member's id
	buildId, buildIdOk := metaData["buildId"].(uuid.UUID)                           // get buildId

	if !skillIdOk || !memberIdOk || !buildIdOk {
		t.Fatalf("Failed to retrieve metadata: skillIdOk=%v, memberIdOk=%v, buildIdOk=%v", skillIdOk, memberIdOk, buildIdOk)
	}

	fmt.Printf("\n@TEST: mainSkillId: %s\nmemberId: %s\nbuildId: %s\n\n", mainSkillId, memberId, buildId)

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
