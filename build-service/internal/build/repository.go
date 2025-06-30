package build

import (
	"fmt"

	// "github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/models"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/types"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllBuilds(
	pageNo int,
	pageSize int,
	sortOrder string,
	sortBy string,
	search string,
	skillId uuid.UUID,
	minRating *int,
	ratingCategory types.RatingCategory) (*[]BuildListQuery, error) {

	var builds []BuildListQuery

	// allowed columns and sort directions
	validSortColumns := map[string]bool{
		"created_at":           true,
		"avg_bossing_rating":   true,
		"avg_endgame_rating":   true,
		"avg_fun_rating":       true,
		"avg_creative_rating":  true,
		"avg_speedfarm_rating": true,
		"main_skill_id":        true,
	}

	validSortDirections := map[string]bool{
		"ASC":  true,
		"DESC": true,
	}

	// validate 'sortBy' and 'sortOrder' and set defaults
	if !validSortColumns[sortBy] {
		sortBy = "created_at"
	}
	if !validSortDirections[sortOrder] {
		sortOrder = "ASC"
	}

	query := `
		SELECT
			builds.id as id,
			title,
			builds.description as description,
			ascendancies.name as ascendancy_name,
			classes.name as class_name,
			skills.name as main_skill_name,
			avg_end_game_rating,
			avg_fun_rating,
			avg_creative_rating,
			avg_speed_farm_rating,
			avg_bossing_rating,
			views,
			status,
			builds.created_at as created_at
		FROM builds
		JOIN classes ON classes.id = builds.class_id
	  JOIN skills ON skills.id = builds.main_skill_id
	  LEFT JOIN ascendancies ON ascendancies.id = builds.ascendancy_id
	`

	// arguments for final query execution
	var queryArgs []interface{}

	// by default status is published when searching for community builds
	query += fmt.Sprintf("\nWHERE status = $1")
	queryArgs = append(queryArgs, types.IsPublished)

	// keyword search filter
	if search != "" {
		searchQuery := "%" + search + "%"
		queryArgs = append(queryArgs, searchQuery)
		query += fmt.Sprintf("\nAND title LIKE $1")
	}

	// skill filter
	if skillId != uuid.Nil {
		queryArgs = append(queryArgs, skillId)
		query += fmt.Sprintf("\nAND main_skill_id = $%d", len(queryArgs))
	}

	// TODO: WIP - rating filter
	// if minRating != nil {
	// 	queryArgs = append(queryArgs, minRating)
	// 	query += fmt.Sprintf("\nAND main_skill_id = $%d", len(queryArgs))
	// }

	// construct pagination and sorting
	queryArgs = append(queryArgs, pageSize, (pageNo-1)*pageSize)

	query += fmt.Sprintf(`
		ORDER BY %s %s
		LIMIT $%d
		OFFSET $%d
	`, sortBy, sortOrder,
		len(queryArgs)-1, // second last arg is LIMIT
		len(queryArgs),   // last arg is OFFSET
	)

	fmt.Printf("\n\nFinal Query: %s\n\n", query)
	fmt.Printf("\n\nFinal QueryArgs: %+v\n\n", queryArgs)

	err := r.db.Select(&builds, query, queryArgs...)

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	buildList := make([]BuildListQuery, len(builds))

	fmt.Println("Builds after query:", builds)

	for index, build := range builds {
		buildList[index] = BuildListQuery{
			Id:                 build.Id,
			Title:              build.Title,
			Description:        build.Description,
			Class:              build.Class,
			Ascendancy:         build.Ascendancy,
			MainSkillName:      build.MainSkillName,
			AvgEndGameRating:   build.AvgEndGameRating,
			AvgFunRating:       build.AvgFunRating,
			AvgCreativeRating:  build.AvgCreativeRating,
			AvgSpeedFarmRating: build.AvgSpeedFarmRating,
			AvgBossingRating:   build.AvgBossingRating,
			Views:              build.Views,
			Status:             build.Status,
			CreatedAt:          build.CreatedAt,
		}
	}

	fmt.Println("all community buildList:", buildList)

	return &buildList, nil
}

func (r *repository) CreateBuild(memberId uuid.UUID, createBuildRequest CreateBuildRequest) (*uuid.UUID, error) {

	baseQuery := `
	INSERT INTO 
		builds(member_id, main_skill_id, class_id, title, description
	`

	if createBuildRequest.AscendancyId != uuid.Nil {
		baseQuery += ", ascendancy_id"
	}

	endQuery := `
	)
	VALUES($1, $2, $3, $4, $5`

	finalQuery := ""

	if createBuildRequest.AscendancyId != uuid.Nil {
		finalQuery = finalQuery + `, $6)
			RETURNING id
			`
	} else {
		finalQuery = finalQuery + `)
		RETURNING id
			`
	}

	query := baseQuery + endQuery + finalQuery

	fmt.Printf("\nFinal query: %+v\n\n", query)

	var buildId uuid.UUID

	var err error
	if createBuildRequest.AscendancyId != uuid.Nil {
		err = r.db.QueryRowx(query, memberId, createBuildRequest.SkillId, createBuildRequest.ClassId, createBuildRequest.Title, createBuildRequest.Description, createBuildRequest.AscendancyId).Scan(&buildId)
	} else {
		err = r.db.QueryRowx(query, memberId, createBuildRequest.SkillId, createBuildRequest.ClassId, createBuildRequest.Title, createBuildRequest.Description).Scan(&buildId)
	}

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	return &buildId, nil
}

/**
* Queries for a corresponding builds tags.
**/
func (r *repository) GetBuildTagsForMemberById(buildId uuid.UUID) (*[]models.Tag, error) {

	var tags []models.Tag

	query := `
	SELECT
		tags.id as id,
		tags.created_at as created_at,
		tags.name as name
	FROM tags
	JOIN build_tags ON build_tags.tag_id = tags.id
	JOIN builds ON builds.id = build_tags.build_id
	WHERE builds.id = $1
	`

	err := r.db.Select(&tags, query, buildId)

	if err != nil {
		fmt.Println("Errored when querying tags.")
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	return &tags, nil
}

func (r *repository) GetBuildsByMemberId(memberId uuid.UUID) (*[]BuildListQuery, error) {
	var builds []BuildListQuery

	query := `
	SELECT
		builds.id as id,
		title,
		builds.description as description,
		ascendancies.name as ascendancy_name,
		classes.name as class_name,
		skills.name as main_skill_name,
		avg_end_game_rating,
		avg_fun_rating,
		avg_creative_rating,
		avg_speed_farm_rating,
		avg_bossing_rating,
		views,
		status,
		builds.created_at as created_at
 	FROM builds
  JOIN classes ON classes.id = builds.class_id
  JOIN skills ON skills.id = builds.main_skill_id
  LEFT JOIN ascendancies ON ascendancies.id = builds.ascendancy_id

	WHERE member_id = $1
	`

	err := r.db.Select(&builds, query, memberId)

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}
	return &builds, nil
}

/**
* Community Version for public viewing when the build is published.
*
* Getting all information related with builds via joins of
* join table build skills, builds, and skills, in order to reduce DB
* round trips.
**/
func (r *repository) GetBuildInfo(buildId uuid.UUID) (*BuildInfoResponse, error) {
	var buildInfoRows []BuildInfoRows

	query := `
	SELECT
		builds.id as id,
		builds.title as title,
		builds.description as description,
		build_skill_links.id as skill_link_id,
		build_skill_links.name as skill_link_name,
		build_skill_links.is_main as skill_link_is_main,
		skills.id as skill_id,
		skills.name as skill_name,
		skills.type as skill_type
	FROM builds
	JOIN build_skill_links ON build_skill_links.build_id = builds.id
	JOIN build_skill_link_skills ON build_skill_link_skills.build_skill_link_id = build_skill_links.id
	JOIN skills ON skills.id = build_skill_link_skills.skill_id
	WHERE builds.id = $1
	ORDER BY build_skill_links.id
	`

	fmt.Println("buildId", buildId)
	err := r.db.Select(&buildInfoRows, query, buildId)

	fmt.Println("buildInfoRows", buildInfoRows)

	if err != nil {
		fmt.Printf("Error when querying for build info: %s\n", err)
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	var buildItemRows []BuildItemSetResponse
	itemQuery := `
	SELECT
		builds.id AS build_id,
		build_item_sets.id AS set_id,
		build_item_set_items.slot AS set_slot,
		COALESCE(items.unique_item, false) AS unique_item,
		COALESCE(items.id, null) AS item_id,
		COALESCE(items.name, '') AS name,
		COALESCE(items.description, '') AS description,
		COALESCE(items.required_level, null) AS required_level,
		COALESCE(items.required_strength, null) AS required_strength,
		COALESCE(items.required_dexterity, null) AS required_dexterity,
		COALESCE(items.required_intelligence, null) AS required_intelligence,
		COALESCE(items.damage, null) AS damage,
		COALESCE(items.aps, null) AS aps,
		COALESCE(items.crit, null) AS crit,
		COALESCE(items.pdps, null) AS pdps,
		COALESCE(items.edps, null) AS edps,
		COALESCE(items.dps, null) AS dps,
		COALESCE(items.additional, null) AS additional,
		COALESCE(items.stats, null) AS stats,
		COALESCE(items.implicit, null) AS implicit,
		COALESCE(items.slot, '') AS slot,
		COALESCE(items.armour, null) AS armour,
		COALESCE(items.energy_shield, null) AS energy_shield,
		COALESCE(items.evasion, null) AS evasion,
		COALESCE(items.block, null) AS block,
		COALESCE(items.ward, null) AS ward
	FROM builds
	JOIN build_item_sets ON build_item_sets.build_id = builds.id
	JOIN build_item_set_items ON build_item_set_items.build_item_set_id = build_item_sets.id
	LEFT JOIN items ON items.id = build_item_set_items.item_id
	WHERE builds.id = $1
	`
	err = r.db.Select(&buildItemRows, itemQuery, buildId)

	if err != nil {
		fmt.Printf("Error when querying for build items: %s\n", err)
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	fmt.Println("buildItemRows", &buildItemRows)
	if len(buildInfoRows) == 0 {
		fmt.Println("No builds queried.")

		// no builds queried with skills or item joins
		return nil, nil
	}

	// create the base of the response
	result := BuildInfoResponse{
		ID:          buildInfoRows[0].ID,
		Title:       buildInfoRows[0].Title,
		Description: buildInfoRows[0].Description,
	}

	var skillRows []models.SkillRow

	for _, buildInfoRow := range buildInfoRows {
		skillRows = append(skillRows, models.SkillRow{
			SkillLinkID:     buildInfoRow.SkillLinkID,
			SkillLinkName:   buildInfoRow.SkillLinkName,
			SkillLinkIsMain: buildInfoRow.SkillLinkIsMain,
			SkillID:         buildInfoRow.SkillID,
			SkillType:       buildInfoRow.SkillType,
		})

	}

	skills := r.GetAndFormSkillLinks(skillRows)

	result.Skills = &skills

	result.Sets = buildItemRows

	return &result, nil
}

/**
* Retrieves and organizes all skills and skill links.
**/

func (r *repository) GetAndFormSkillLinks(skillData []models.SkillRow) SkillGroupResponse {
	var mainSkillLink SkillLinkResponse          // store primary skills
	var additionalSkillLinks []SkillLinkResponse // stores additional skills

	// group up all skill information
	for _, row := range skillData {

		// --- grouping primary skills ---

		// identify the "main skilllink" with the "skill_link_is_main" field
		if row.SkillLinkIsMain {
			mainSkillLink.SkillLinkName = row.SkillLinkName

			// match start of skill by active skill - match after casting
			if types.SkillType(row.SkillType) == types.Active {

				mainSkillLink.Skill = models.Skill{
					Id:   row.SkillID,
					Name: row.SkillName,
					Type: row.SkillType,
				}

			} else {
				// else its a support skill link
				mainSkillLink.Links = append(mainSkillLink.Links, models.Skill{
					Id:   row.SkillID,
					Name: row.SkillName,
					Type: row.SkillType,
				})
			}
		} else {
			// else we construct the secondary skills

			// --- grouping secondary skills ---

			// find existing skillLink via SkillLinkName
			skillLinkExists := false
			var existingSkillLink *SkillLinkResponse

			for index := range additionalSkillLinks {
				if additionalSkillLinks[index].SkillLinkName == row.SkillLinkName {
					skillLinkExists = true
					// save reference to original skill link slice
					existingSkillLink = &additionalSkillLinks[index]
					break
				}
			}

			// update existing link
			if skillLinkExists {
				// starting link skill
				if types.SkillType(row.SkillType) == types.Active {

					existingSkillLink.Skill = models.Skill{
						Id:   row.SkillID,
						Name: row.SkillName,
						Type: row.SkillType,
					}

				} else {
					existingSkillLink.Links = append(existingSkillLink.Links, models.Skill{
						Id:   row.SkillID,
						Name: row.SkillName,
						Type: row.SkillType,
					})
				}

			} else {
				// creating new link

				var newAdditionalSkillLink SkillLinkResponse

				// create new skill link name and the first skill
				newAdditionalSkillLink.SkillLinkName = row.SkillLinkName

				// starting link skill
				if types.SkillType(row.SkillType) == types.Active {
					newAdditionalSkillLink.Skill = models.Skill{
						Id:   row.SkillID,
						Name: row.SkillName,
						Type: row.SkillType,
					}

				} else {
					// supporting skill link
					newAdditionalSkillLink.Links = append(newAdditionalSkillLink.Links, models.Skill{
						Id:   row.SkillID,
						Name: row.SkillName,
						Type: row.SkillType,
					})
				}

				// add it to slice of additionalSkillLinks
				additionalSkillLinks = append(additionalSkillLinks, newAdditionalSkillLink)
			}

		}

	}

	// wrap them for response
	skills := SkillGroupResponse{
		MainSkillLinks:   mainSkillLink,
		AdditionalSkills: additionalSkillLinks,
	}

	return skills
}
