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
	ratingCategory types.RatingCategory) ([]BuildListQuery, error) {

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

	return buildList, nil
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
