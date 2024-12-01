package build_test

import (
	"fmt"
	"testing"

	"github.com/darkphotonKN/community-builds/internal/build"
	"github.com/darkphotonKN/community-builds/internal/skill"
	testsuite "github.com/darkphotonKN/community-builds/internal/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type SkillLinkRow struct {
	Name      string    `db:"name"`
	IsMain    bool      `db:"is_main"`
	SkillID   uuid.UUID `db:"skill_id"`
	SkillName string    `db:"skill_name"`
}

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

	// query skill links
	var skillLinkRows []SkillLinkRow

	err := ts.DB.Select(&skillLinkRows, `
		SELECT 
			build_skill_links.name AS name,
			build_skill_links.is_main AS is_main,
			skills.id AS skill_id,
			skills.name as skill_name
		FROM build_skill_links
		JOIN build_skill_link_skills ON build_skill_link_skills.build_skill_link_id = build_skill_links.id
		JOIN skills ON skills.id = build_skill_link_skills.skill_id
		WHERE build_skill_links.build_id = $1
	`, buildId)

	if err != nil {
		t.Fatalf("Could not find build that was supposed to be created. Error: %s\n", err)
	}

	fmt.Printf("\nskillLinkRows: result: %+v\n", skillLinkRows)

	// test result count is expected

	// 1 main skill + main skill links + additional skill links
	expectedCount := 1 // default 1 from main skill

	// count main skill links
	for range payload.MainSkillLinks.Links {
		expectedCount++
	}

	// count additional skill links
	for _, links := range payload.AdditionalSkills {
		// add one for each additional skill
		expectedCount++
		for range links.Links {
			// add one for each link
			expectedCount++
		}
	}

	fmt.Printf("\ntargetCount: %d\n", expectedCount)
	fmt.Printf("\nlen(skillLinkRows) %d\n", len(skillLinkRows))

	assert.Len(t, skillLinkRows, expectedCount, fmt.Sprintf("Expected %d skill links, got %d", expectedCount, len(skillLinkRows)))

	// test main skill is earthquake and links are correct

	var isMainSkillRes []SkillLinkRow

	// get all the skills that were marked as "is_main"
	for _, skill := range skillLinkRows {
		if skill.IsMain {
			isMainSkillRes = append(isMainSkillRes, skill)
		}
	}

	var expectedMainSkills []struct {
		ID   uuid.UUID
		Name string
	}

	// get all the main skills of payload
	expectedMainSkills = append(expectedMainSkills, struct {
		ID   uuid.UUID
		Name string
	}{
		testSkill1.ID,
		testSkill1.Name,
	})

	expectedMainSkills = append(expectedMainSkills, struct {
		ID   uuid.UUID
		Name string
	}{
		testSkill2.ID,
		testSkill2.Name,
	})

	expectedMainSkills = append(expectedMainSkills, struct {
		ID   uuid.UUID
		Name string
	}{
		testSkill3.ID,
		testSkill3.Name,
	})

	// test against results
	mainSkillMatches := 0

	for _, mainSkill := range isMainSkillRes {
		// check if it exists in the isMainSkills
		for _, expectedMainSkill := range expectedMainSkills {
			if mainSkill.SkillID == expectedMainSkill.ID {
				// name must be the same
				assert.Equal(t, expectedMainSkill.Name, mainSkill.SkillName, fmt.Sprintf("expected main skill name: %s but got: %s\n", expectedMainSkill.Name, mainSkill.SkillName))

				// count main skill
				mainSkillMatches++
			}
		}
	}

	// check there were 3 main skills created
	assert.Equal(t, 3, mainSkillMatches, "Not enough main skills were created.")

	// test additional skills

	var isAdditionalSkill []SkillLinkRow

	// get all the skills that were secondary, NOT marked as "is_main"
	for _, skill := range skillLinkRows {
		if !skill.IsMain {
			isAdditionalSkill = append(isAdditionalSkill, skill)
		}
	}

	var expectedAdditionalSkills []struct {
		ID   uuid.UUID
		Name string
	}

	// get all the main skills of payload
	expectedAdditionalSkills = append(expectedAdditionalSkills, struct {
		ID   uuid.UUID
		Name string
	}{
		testSkill4.ID,
		testSkill4.Name,
	})

	expectedAdditionalSkills = append(expectedAdditionalSkills, struct {
		ID   uuid.UUID
		Name string
	}{
		testSkill5.ID,
		testSkill5.Name,
	})

	// test against results
	additionalSkillMatches := 0

	for _, additionalSkill := range isAdditionalSkill {
		// check if it exists in the isMainSkills
		for _, expectedAdditionalSkill := range expectedAdditionalSkills {
			if additionalSkill.SkillID == expectedAdditionalSkill.ID {
				// name must be the same
				assert.Equal(t, expectedAdditionalSkill.Name, additionalSkill.SkillName, fmt.Sprintf("expected additional skill name: %s but got: %s\n", expectedAdditionalSkill.Name, additionalSkill.SkillName))

				// count main skill
				additionalSkillMatches++
			}
		}
	}

	assert.Equal(t, 2, additionalSkillMatches)
}
