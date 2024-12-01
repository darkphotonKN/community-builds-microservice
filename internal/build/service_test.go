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

	ts := testsuite.NewTestSuite()

	// -- retrieve stored metaData --
	metaData := ts.GetMetaData("TestAddSkillLinksToBuildService_Success") // global metadata
	mainSkillId, skillIdOk := metaData["mainSkillId"].(uuid.UUID)         // get main skill id
	memberId, memberIdOk := metaData["memberId"].(uuid.UUID)              // get a member's id
	buildId, buildIdOk := metaData["buildId"].(uuid.UUID)                 // get buildId
	skills, skillsOk := metaData["skills"].(testsuite.TestSkills)         // get skill information

	if !skillIdOk || !memberIdOk || !buildIdOk || !skillsOk {
		t.Fatalf("Failed to retrieve metadata: skillIdOk=%v, memberIdOk=%v, buildIdOk=%v, skillsOk=%v", skillIdOk, memberIdOk, buildIdOk, skillsOk)
	}

	fmt.Printf("\n@TEST: mainSkillId: %s\nmemberId: %s\nbuildId: %s\nskills: %+v\n\n", mainSkillId, memberId, buildId, skills)

	// --- test service methods ---

	// -- DI setup --
	skillRepo := skill.NewSkillRepository(ts.DB)
	skillService := skill.NewSkillService(skillRepo)

	buildsRepo := build.NewBuildRepository(ts.DB)
	buildsService := build.NewBuildService(buildsRepo, skillService)

	fmt.Println(buildsService)

	// -- tests --

	// find and insert all test skills
	var testSkill1 testsuite.TestSkill
	var testSkill2 testsuite.TestSkill
	var testSkill3 testsuite.TestSkill
	var testSkill4 testsuite.TestSkill
	var testSkill5 testsuite.TestSkill

	for _, skill := range skills {
		switch skill.Name {
		case "Earthquake":
			testSkill1 = skill
		case "Concentrated Effect":
			testSkill2 = skill
		case "Increased Area of Effect":
			testSkill3 = skill
		case "Lightning Strike":
			testSkill4 = skill
		case "Multistrike":
			testSkill5 = skill
		}
	}

	fmt.Println("testSkill1:", testSkill1)

	// mock payload
	payload := build.AddSkillsToBuildRequest{
		MainSkillLinks: build.SkillLinks{
			SkillLinkName: "Earthquake",
			Skill:         mainSkillId,
			Links:         []uuid.UUID{testSkill2.ID, testSkill3.ID},
		},
		AdditionalSkills: []build.SkillLinks{
			{
				SkillLinkName: "Lighting Strike Multistrike",
				Skill:         testSkill4.ID,
				Links:         []uuid.UUID{testSkill5.ID},
			},
		},
	}

	// test method
	buildsService.AddSkillLinksToBuildService(memberId, buildId, payload)

	// assert skill links count
	// expectedSkillLinksCount := len(payload.AdditionalSkills) + 1 // +1 for MainSkillLinks
	//
	//	if len(skillLinks) != expectedSkillLinksCount {
	//		t.Fatalf("Expected %d skill links, got %d", expectedSkillLinksCount, len(skillLinks))
	//	}
}
