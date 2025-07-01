package build

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/skill"
	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/tag"
	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/utils/dbutils"
	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/utils/errorutils"
	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/build"
	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	db           *sqlx.DB
	repo         Repository
	publishCh    *amqp.Channel
	skillService skill.Service
	tagService   tag.Service
}

type Repository interface {
	CreateBuild(memberId uuid.UUID, createBuildRequest CreateBuildRequest) (*uuid.UUID, error)
	GetBuildsForMember(id uuid.UUID) (*[]BuildListQuery, error)
	GetAllBuilds(pageNo int, pageSize int, sortOrder string, sortBy string, search string, skillId uuid.UUID, minRating *int, ratingCategory types.RatingCategory) (*[]BuildListQuery, error)
	GetBuildInfo(buildId uuid.UUID) (*BuildInfoResponse, error)
	GetAndFormSkillLinks(skillData []models.SkillRow) SkillGroupResponse
	CreateBuildTags(buildId uuid.UUID, tagIds []uuid.UUID) error
	GetBuildForMemberById(memberId uuid.UUID, buildId uuid.UUID) (*models.Build, error)
	CreateItemToSetTx(tx *sqlx.Tx, buildItemSetId uuid.UUID, slot string, itemId interface{}) error
	CreateBuildItemSetTx(tx *sqlx.Tx, buildId uuid.UUID) (uuid.UUID, error)
	GetBuildTagsForMemberById(buildId uuid.UUID) (*[]models.Tag, error)
	GetBuildAscendancyById(buildId uuid.UUID) (*string, error)
	GetBuildClassById(buildId uuid.UUID) (*string, error)
	GetBuildItemSetById(buildId uuid.UUID) ([]BuildItemSetResponse, error)
	GetBasicBuildInfoByIdForMember(buildId uuid.UUID, MemberId uuid.UUID) (*BasicBuildInfoResponse, error)
	UpdateBuildByIdForMember(id uuid.UUID, memberId uuid.UUID, updateFields models.Build) error
	UpdateBuild(memberId uuid.UUID, buildId uuid.UUID, request UpdateBuildRequest) error
	CreateBuildSkillLinkTx(tx *sqlx.Tx, buildId uuid.UUID, name string, isMain bool) (uuid.UUID, error)
	AddSkillToLinkTx(tx *sqlx.Tx, buildSkillLinkId uuid.UUID, skillId uuid.UUID) error
	GetBuildItemSetIdTx(tx *sqlx.Tx, buildId uuid.UUID) (uuid.UUID, error)
	UpdateItemToSetTx(tx *sqlx.Tx, buildItemSetId uuid.UUID, slot string, itemId interface{}) error
	DeleteBuildByIdForMember(memberId uuid.UUID, buildId uuid.UUID) error
}

func NewService(db *sqlx.DB, repo Repository, publishCh *amqp.Channel, skillService skill.Service, tagService tag.Service) Service {
	return &service{db: db, repo: repo, publishCh: publishCh, skillService: skillService, tagService: tagService}
}

/**
* Get list of builds available to a member.
**/

func (s *service) GetBuildsForMember(ctx context.Context, req *pb.GetBuildsForMemberRequest) (*pb.GetBuildsForMemberResponse, error) {
	memberId, err := uuid.Parse(req.MemberId)

	if err != nil {
		return nil, err
	}

	baseBuilds, err := s.repo.GetBuildsForMember(memberId)

	if err != nil {
		return nil, err
	}

	pbBuilds := make([]*pb.BuildList, len(*baseBuilds))

	// query and add each builds tag's
	for index, build := range *baseBuilds {
		tags, err := s.GetBuildTagsForMemberById(build.Id)

		// stop query pre-maturely if errored on query
		if err != nil {
			return nil, err
		}

		pbTags := make([]*pb.Tag, len(*tags))
		for tagIndex, tag := range *tags {
			pbTags[tagIndex] = &pb.Tag{
				Id:        tag.ID.String(),
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt.String(),
				UpdatedAt: tag.UpdatedAt.String(),
			}
		}
		pbBuilds[index] = &pb.BuildList{
			Id:                 build.Id.String(),
			Title:              build.Title,
			Description:        build.Description,
			Class:              build.Class,
			Ascendancy:         *build.Ascendancy,
			MainSkillName:      build.MainSkillName,
			AvgEndGameRating:   *build.AvgEndGameRating,
			AvgBossingRating:   *build.AvgBossingRating,
			AvgCreativeRating:  *build.AvgCreativeRating,
			AvgFunRating:       *build.AvgFunRating,
			AvgSpeedFarmRating: *build.AvgSpeedFarmRating,
			Tags:               pbTags,
			Views:              int32(build.Views),
			Status:             int32(build.Status),
			CreatedAt:          build.CreatedAt,
		}
	}

	return &pb.GetBuildsForMemberResponse{
		Builds: pbBuilds,
	}, nil
}

// func (s *service) GetBuildsForMember(ctx context.Context, req *pb.GetBuildsForMemberRequest) (*pb.GetBuildsForMemberResponse, error) {
// 	MemberId, err := uuid.Parse(req.MemberId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	baseBuilds, err := s.repo.GetBuildsForMember(MemberId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	buildListResponse := make([]BuildListResponse, len(*baseBuilds))

// 	// query and add each builds tag's
// 	for index, build := range *baseBuilds {
// 		tags, err := s.tagService.GetBuildTagsForMemberById(build.Id)

// 		// stop query pre-maturely if errored on query
// 		if err != nil {
// 			return nil, err
// 		}

// 		buildListResponse[index] = BuildListResponse{
// 			Id:                 build.Id,
// 			Title:              build.Title,
// 			Description:        build.Description,
// 			Class:              build.Class,
// 			Ascendancy:         build.Ascendancy,
// 			MainSkillName:      build.MainSkillName,
// 			AvgEndGameRating:   build.AvgEndGameRating,
// 			AvgBossingRating:   build.AvgBossingRating,
// 			AvgCreativeRating:  build.AvgCreativeRating,
// 			AvgFunRating:       build.AvgFunRating,
// 			AvgSpeedFarmRating: build.AvgSpeedFarmRating,
// 			Tags:               *tags,
// 			Views:              build.Views,
// 			Status:             build.Status,
// 			CreatedAt:          build.CreatedAt,
// 		}
// 	}

// 	return &buildListResponse, nil
// }

const (
	maxBuildCount = 10
)

/**
* Gets list of public community builds.
**/

func (s *service) GetCommunityBuilds(ctx context.Context, req *pb.GetCommunityBuildsRequest) (*pb.GetCommunityBuildsResponse, error) {

	pageNo := int(req.PageNo)
	pageSize := int(req.PageSize)
	sortOrder := string(req.SortOrder)
	sortBy := string(req.SortBy)
	search := string(req.Search)
	skillId, err := uuid.Parse(req.SkillId)
	if err != nil {
		return nil, err
	}
	minRating := int(req.MinRating)
	ratingCategory := types.RatingCategory(req.RatingCategory)

	baseBuilds, err := s.repo.GetAllBuilds(pageNo, pageSize, sortOrder, sortBy, search, skillId, &minRating, ratingCategory)

	if err != nil {
		return nil, err
	}

	pbBuilds := make([]*pb.BuildList, len(*baseBuilds))
	// query and add each builds tag's
	for index, build := range *baseBuilds {
		tags, err := s.GetBuildTagsForMemberById(build.Id)

		// stop query pre-maturely if errored on query
		if err != nil {
			return nil, err
		}

		idStr := build.Id.String()
		pbTags := make([]*pb.Tag, len(*tags))
		for tagIndex, tag := range *tags {
			pbTags[tagIndex] = &pb.Tag{
				Id:        tag.ID.String(),
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt.String(),
				UpdatedAt: tag.UpdatedAt.String(),
			}
		}
		pbBuilds[index] = &pb.BuildList{
			Id:                 idStr,
			Title:              build.Title,
			Description:        build.Description,
			Class:              build.Class,
			Ascendancy:         *build.Ascendancy,
			MainSkillName:      build.MainSkillName,
			AvgEndGameRating:   *build.AvgEndGameRating,
			AvgBossingRating:   *build.AvgBossingRating,
			AvgCreativeRating:  *build.AvgCreativeRating,
			AvgFunRating:       *build.AvgFunRating,
			AvgSpeedFarmRating: *build.AvgSpeedFarmRating,
			Tags:               pbTags,
			Views:              int32(build.Views),
			Status:             int32(build.Status),
			CreatedAt:          build.CreatedAt,
		}
	}

	return &pb.GetCommunityBuildsResponse{
		Builds: pbBuilds,
	}, nil
}

func (s *service) CreateBuild(ctx context.Context, req *pb.CreateBuildRequest) (*pb.CreateBuildResponse, error) {
	fmt.Println("req", req)
	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}

	skillId, err := uuid.Parse(req.SkillId)
	if err != nil {
		return nil, err
	}
	// confirm skill exists
	_, err = s.skillService.GetSkillById(skillId)

	if err != nil {
		return nil, fmt.Errorf("main skill id could not be found when attempting to create build for it.")
	}

	// check how many builds a user has, prevent creation if over the limit
	buildReq, err := s.GetBuildsForMember(ctx, &pb.GetBuildsForMemberRequest{MemberId: req.MemberId})

	if err != nil {
		return nil, err
	}

	if len(buildReq.Builds) > maxBuildCount {
		return nil, fmt.Errorf("Number of builds allowed have reached maximum capacity.")
	}

	SkillId, err := uuid.Parse(req.SkillId)
	if err != nil {
		return nil, err
	}

	ClassId, err := uuid.Parse(req.ClassId)
	if err != nil {
		return nil, err
	}

	AscendancyId, err := uuid.Parse(req.AscendancyId)
	if err != nil {
		return nil, err
	}

	TagIds := make([]uuid.UUID, len(req.TagIds))
	for _, tagId := range req.TagIds {
		TagIds = append(TagIds, uuid.Must(uuid.Parse(tagId)))
	}

	// format to fit model for db tags
	createBuild := &CreateBuildRequest{
		SkillId:      SkillId,
		TagIds:       TagIds,          // this is a list of tag ids
		Title:        req.Title,       // Title is the name of the build
		Description:  req.Description, // Description of the build
		ClassId:      ClassId,         // Class ID of the build
		AscendancyId: AscendancyId,    // Ascendancy ID of the build
	}

	// create build with this skill and for this member
	buildId, err := s.repo.CreateBuild(memberId, *createBuild)
	fmt.Println("buildId:", buildId)
	if err != nil {
		return nil, err
	}

	// create build tags
	var tagIds = make([]uuid.UUID, len(req.TagIds))
	for _, tag := range req.TagIds {
		parseTag, err := uuid.Parse(tag)
		if err == nil {
			tagIds = append(tagIds, parseTag)
		}
	}
	err = s.repo.CreateBuildTags(*buildId, tagIds)

	if err != nil {
		return nil, err
	}

	// create build default set
	err = s.CreateDefaultItemSetsToBuild(memberId, *buildId)

	if err != nil {
		return nil, err
	}

	// return nil

	// publish rabbit mq message after succesfuly creating
	marshalledBuild, err := json.Marshal(createBuild)

	if err != nil {
		return nil, err
	}

	err = s.publishCh.PublishWithContext(
		ctx,
		commonconstants.BuildCreatedEvent,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshalledBuild,
			// persist message
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		return nil, err
	}

	return &pb.CreateBuildResponse{}, nil
}

func (s *service) GetBuildInfo(ctx context.Context, req *pb.GetBuildInfoRequest) (*pb.GetBuildInfoResponse, error) {

	// return all join information (base, class, ascendancy, skills and items)
	buildId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	info, err := s.repo.GetBuildInfo(buildId)

	pbTags := make([]*pb.Tag, len(info.Tags))
	for tagIndex, tag := range info.Tags {
		pbTags[tagIndex] = &pb.Tag{
			Id:        tag.ID.String(),
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt.String(),
			UpdatedAt: tag.UpdatedAt.String(),
		}
	}

	pbLinks := make([]*pb.Skill, len(info.Skills.MainSkillLinks.Links))
	for skillIndex, skill := range info.Skills.MainSkillLinks.Links {
		pbLinks[skillIndex] = &pb.Skill{
			Id:        skill.Id.String(),
			Name:      skill.Name,
			Type:      skill.Type,
			CreatedAt: skill.CreatedAt.String(),
			UpdatedAt: skill.UpdatedAt.String(),
		}
	}
	pbSkill := &pb.Skill{
		Id:        info.Skills.MainSkillLinks.Skill.Id.String(),
		Name:      info.Skills.MainSkillLinks.Skill.Name,
		Type:      info.Skills.MainSkillLinks.Skill.Type,
		CreatedAt: info.Skills.MainSkillLinks.Skill.CreatedAt.String(),
		UpdatedAt: info.Skills.MainSkillLinks.Skill.UpdatedAt.String(),
	}

	pbSkills := &pb.SkillGroupResponse{
		MainSkillLinks: &pb.SkillLinkResponse{
			SkillLinkName: info.Skills.MainSkillLinks.SkillLinkName,
			Skill:         pbSkill,
			Links:         pbLinks,
		},
	}

	pbSets := make([]*pb.BuildItemSetResponse, len(info.Sets))
	for SetIndex, Set := range info.Sets {
		pbSets[SetIndex] = &pb.BuildItemSetResponse{
			BuildId:     Set.BuildId.String(),
			SetId:       Set.SetId.String(),
			ItemId:      Set.ItemId.String(),
			ImageUrl:    Set.ImageUrl,
			Category:    Set.Category,
			Class:       Set.Class,
			Name:        Set.Name,
			Type:        Set.Type,
			Description: Set.Description,
			UniqueItem:  Set.UniqueItem,
			Slot:        Set.Slot,
			// todo 補齊剩下的選填

			RequiredLevel:        Set.RequiredLevel,
			RequiredStrength:     Set.RequiredStrength,
			RequiredDexterity:    Set.RequiredDexterity,
			RequiredIntelligence: Set.RequiredIntelligence,
			Armour:               Set.Armour,
			EnergyShield:         Set.EnergyShield,
			Evasion:              Set.Evasion,
			Block:                Set.Block,
			Ward:                 Set.Ward,

			Damage: Set.Damage,
			APS:    Set.APS,
			Crit:   Set.Crit,
			PDPS:   Set.PDPS,
			EDPS:   Set.EDPS,
			DPS:    Set.DPS,

			Life:     Set.Life,
			Mana:     Set.Mana,
			Duration: Set.Duration,
			Usage:    Set.Usage,
			Capacity: Set.Capacity,

			Additional: Set.Additional,
			Stats:      Set.Stats,
			Implicit:   *Set.Implicit,
		}
	}

	res := &pb.GetBuildInfoResponse{
		Id:          info.Id.String(),
		Title:       info.Title,
		Description: info.Description,
		Class:       info.Class,
		Ascendancy:  *info.Ascendancy,
		Tags:        pbTags,
		Skills:      pbSkills,
		Sets:        pbSets,
	}

	return res, nil
}

/**
* Create default set,
* rolling back on error.
**/
func (s *service) CreateDefaultItemSetsToBuild(memberId uuid.UUID, buildId uuid.UUID) error {

	return dbutils.ExecTx(s.db, func(tx *sqlx.Tx) error {
		// get build and check if it exists
		_, err := s.GetBuildForMemberById(memberId, buildId)

		if err != nil {
			return err
		}

		itemSetId, err := s.repo.CreateBuildItemSetTx(tx, buildId)

		if err != nil {
			return err
		}

		itemSetsMap := map[string]interface{}{
			"weapon":     "",
			"shield":     "",
			"helmet":     "",
			"bodyArmour": "",
			"gloves":     "",
			"belt":       "",
			"boots":      "",
			"amulet":     "",
			"leftRing":   "",
			"rightRing":  "",
		}

		// create item relations under this item set, one item at a time
		for key, value := range itemSetsMap {
			if value == "" {
				value = nil
			}
			err = s.repo.CreateItemToSetTx(tx, itemSetId, key, value)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *service) GetBuildInfoForMember(ctx context.Context, req *pb.GetBuildInfoForMemberRequest) (*pb.GetBuildInfoForMemberResponse, error) {
	// base build
	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}
	buildId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	build, err := s.GetBuildForMemberById(memberId, buildId)

	if err != nil {
		return nil, err
	}

	// tags
	tags, err := s.repo.GetBuildTagsForMemberById(buildId)

	if err != nil {
		return nil, err
	}

	// retrieve build skills

	// get all skill and link info that exist in a build
	skillRows, err := s.skillService.GetSkillsByBuildId(buildId)

	if err != nil {
		return nil, err
	}

	fmt.Printf("\nSkillrows: %+v\n\n", skillRows)

	// group them into the skill group response, if there are returned data
	var skills *SkillGroupResponse

	if len(*skillRows) > 0 {
		formedSkills := s.repo.GetAndFormSkillLinks(*skillRows)
		skills = &formedSkills
	}

	fmt.Printf("\nFormed Skills: %+v\n\n", skills)

	if err != nil {
		return nil, err
	}

	// get class (name only)
	class, err := s.repo.GetBuildClassById(buildId)

	if err != nil {
		return nil, err
	}

	// get ascendancy (name only)
	ascendancy, err := s.repo.GetBuildAscendancyById(buildId)

	if err != nil {
		return nil, err
	}

	// TODO: retrieve build items
	set, err := s.repo.GetBuildItemSetById(buildId)
	if err != nil {
		fmt.Println("set err", err)
		return nil, err
	}
	fmt.Println("set", set)

	pbLinks := make([]*pb.Skill, len(skills.MainSkillLinks.Links))
	for skillIndex, skill := range skills.MainSkillLinks.Links {
		pbLinks[skillIndex] = &pb.Skill{
			Id:        skill.Id.String(),
			Name:      skill.Name,
			Type:      skill.Type,
			CreatedAt: skill.CreatedAt.String(),
			UpdatedAt: skill.UpdatedAt.String(),
		}
	}
	pbSkill := &pb.Skill{
		Id:        skills.MainSkillLinks.Skill.Id.String(),
		Name:      skills.MainSkillLinks.Skill.Name,
		Type:      skills.MainSkillLinks.Skill.Type,
		CreatedAt: skills.MainSkillLinks.Skill.CreatedAt.String(),
		UpdatedAt: skills.MainSkillLinks.Skill.UpdatedAt.String(),
	}

	pbSkills := &pb.SkillGroupResponse{
		MainSkillLinks: &pb.SkillLinkResponse{
			SkillLinkName: skills.MainSkillLinks.SkillLinkName,
			Skill:         pbSkill,
			Links:         pbLinks,
		},
	}

	pbTags := make([]*pb.Tag, len(*tags))
	for tagIndex, tag := range *tags {
		pbTags[tagIndex] = &pb.Tag{
			Id:        tag.ID.String(),
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt.String(),
			UpdatedAt: tag.UpdatedAt.String(),
		}
	}

	pbSets := make([]*pb.BuildItemSetResponse, len(set))
	for SetIndex, Set := range set {
		pbSets[SetIndex] = &pb.BuildItemSetResponse{
			BuildId:     Set.BuildId.String(),
			SetId:       Set.SetId.String(),
			ItemId:      Set.ItemId.String(),
			ImageUrl:    Set.ImageUrl,
			Category:    Set.Category,
			Class:       Set.Class,
			Name:        Set.Name,
			Type:        Set.Type,
			Description: Set.Description,
			UniqueItem:  Set.UniqueItem,
			Slot:        Set.Slot,
			// todo 補齊剩下的選填

			RequiredLevel:        Set.RequiredLevel,
			RequiredStrength:     Set.RequiredStrength,
			RequiredDexterity:    Set.RequiredDexterity,
			RequiredIntelligence: Set.RequiredIntelligence,
			Armour:               Set.Armour,
			EnergyShield:         Set.EnergyShield,
			Evasion:              Set.Evasion,
			Block:                Set.Block,
			Ward:                 Set.Ward,

			Damage: Set.Damage,
			APS:    Set.APS,
			Crit:   Set.Crit,
			PDPS:   Set.PDPS,
			EDPS:   Set.EDPS,
			DPS:    Set.DPS,

			Life:     Set.Life,
			Mana:     Set.Mana,
			Duration: Set.Duration,
			Usage:    Set.Usage,
			Capacity: Set.Capacity,

			Additional: Set.Additional,
			Stats:      Set.Stats,
			Implicit:   *Set.Implicit,
		}
	}

	grpcReq := &pb.GetBuildInfoForMemberResponse{
		Id:          build.ID.String(),
		Title:       build.Title,
		Description: build.Description,
		// TODO: add ascendancy and class
		Class:      *class,
		Ascendancy: *ascendancy,
		Skills:     pbSkills,
		Tags:       pbTags,
		Sets:       pbSets,
	}

	fmt.Printf("Constructed build info: %+v\n", grpcReq)

	return grpcReq, nil

}

/**
* Get a single build for a member by Id, without extra joined information.
**/
func (s *service) GetBuildForMemberById(memberId uuid.UUID, buildId uuid.UUID) (*models.Build, error) {
	return s.repo.GetBuildForMemberById(memberId, buildId)
}

// not grpc
func (s *service) GetBuildTagsForMemberById(buildId uuid.UUID) (*[]models.Tag, error) {
	tags, err := s.repo.GetBuildTagsForMemberById(buildId)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *service) PublishBuild(ctx context.Context, req *pb.PublishBuildRequest) (*pb.PublishBuildResponse, error) {

	// check if build belongs to member
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}
	build, err := s.repo.GetBasicBuildInfoByIdForMember(id, memberId)

	fmt.Printf("Build retrieved for member: %+v", build)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	if build == nil {
		return nil, errors.New(fmt.Sprintf("No build with the id %s exists.\n", id))
	}

	publishBuild := models.Build{
		Status: int(published),
	}
	s.repo.UpdateBuildByIdForMember(id, memberId, publishBuild)

	return &pb.PublishBuildResponse{}, nil
}

func (s *service) UpdateBuild(ctx context.Context, req *pb.UpdateBuildRequest) (*pb.UpdateBuildResponse, error) {
	// TODO: confirm skill exists
	// _, err := s.SkillService.GetSkillByIdService(request.SkillID)

	// if err != nil {
	// 	return fmt.Errorf("main skill id could not be found when attempting to update build for it.")
	// }

	// check build exists before updating
	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}
	buildId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.GetBuildForMemberById(memberId, buildId)

	if err != nil {
		return nil, fmt.Errorf("Can only update an existing build. %s", err)
	}
	skillId, err := uuid.Parse(req.SkillId)
	if err != nil {
		return nil, err
	}
	tags := make([]uuid.UUID, len(req.Tags))
	for _, tag := range req.Tags {

		uuidTag, _ := uuid.Parse(tag)
		tags = append(tags, uuidTag)
	}

	classId, err := uuid.Parse(req.ClassId)
	if err != nil {
		return nil, err
	}
	ascendancyId, err := uuid.Parse(req.AscendancyId)
	if err != nil {
		return nil, err
	}
	request := UpdateBuildRequest{
		SkillId:      &skillId,
		TagIds:       tags,
		Title:        &req.Title,
		Description:  &req.Description,
		ClassId:      &classId,
		AscendancyId: &ascendancyId,
	}

	err = s.repo.UpdateBuild(memberId, buildId, request)

	if err != nil {
		return nil, err
	}

	// TODO: update build tags
	// err = s.Repo.CreateBuildTags(*buildId, request.TagIDs)

	if err != nil {
		return nil, err
	}

	return &pb.UpdateBuildResponse{}, nil
}

func (s *service) AddSkillLinksToBuild(ctx context.Context, req *pb.AddSkillLinksToBuildRequest) (*pb.AddSkillLinksToBuildResponse, error) {
	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}
	buildId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	return nil, dbutils.ExecTx(s.db, func(tx *sqlx.Tx) error {
		// get build and check if it exists

		_, err = s.GetBuildForMemberById(memberId, buildId)

		if err != nil {
			return err
		}

		// -- CREATE SKILL LINKS FOR BUILD --

		// --- primary links ---

		// add main skill group

		mainSkillLinkId, err := s.repo.CreateBuildSkillLinkTx(tx, buildId, req.MainSkillLinks.SkillLinkName, true)

		// add main skill relation to main skill link
		skillId, err := uuid.Parse(req.MainSkillLinks.Skill)
		if err != nil {
			return err
		}
		err = s.repo.AddSkillToLinkTx(tx, mainSkillLinkId, skillId)
		// create item mod

		if err != nil {
			return err
		}

		fmt.Printf("mainSkillLink created, id: %s\n", mainSkillLinkId)

		// create skill relations under this main skill link, one skill at a time
		for _, skillId := range req.MainSkillLinks.Links {
			LinksSkillId, _ := uuid.Parse(skillId)

			err := s.repo.AddSkillToLinkTx(tx, mainSkillLinkId, LinksSkillId)
			if err != nil {
				return err
			}
		}

		// --- other links --
		for _, skillLinks := range req.AdditionalSkills {

			// add secondary skill group
			secondarySkillLinkId, err := s.repo.CreateBuildSkillLinkTx(tx, buildId, skillLinks.SkillLinkName, false)
			Skill, _ := uuid.Parse(skillLinks.Skill)
			// add main skill relation to secondary link
			err = s.repo.AddSkillToLinkTx(tx, secondarySkillLinkId, Skill)
			if err != nil {
				return err
			}

			fmt.Printf("secondarySkillLinkId created, id: %s\n", secondarySkillLinkId)

			// create skill relations under this secondary skill link, one skill at a time
			for _, skillId := range skillLinks.Links {
				LinksSkillId, _ := uuid.Parse(skillId)
				err := s.repo.AddSkillToLinkTx(tx, secondarySkillLinkId, LinksSkillId)

				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

/**
* Adds primary and secondary items and links to an existing build via a transanction,
* rolling back on error.
**/
func (s *service) UpdateItemSetsToBuild(ctx context.Context, req *pb.UpdateItemSetsToBuildRequest) (*pb.UpdateItemSetsToBuildResponse, error) {
	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}
	buildId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	return nil, dbutils.ExecTx(s.db, func(tx *sqlx.Tx) error {
		// get build and check if it exists

		_, err = s.GetBuildForMemberById(memberId, buildId)

		if err != nil {
			return err
		}

		itemSetId, err := s.repo.GetBuildItemSetIdTx(tx, buildId)

		if err != nil {
			return err
		}
		fmt.Println("itemSetId", itemSetId)

		itemSetsMap := map[string]string{
			"weapon":     req.Weapon,
			"shield":     req.Shield,
			"helmet":     req.Helmet,
			"bodyArmour": req.BodyArmour,
			"gloves":     req.Gloves,
			"belt":       req.Belt,
			"boots":      req.Boots,
			"amulet":     req.Amulet,
			"leftRing":   req.LeftRing,
			"rightRing":  req.RightRing,
		}
		// create item relations under this item set, one item at a time
		for key, value := range itemSetsMap {
			// 表示空值
			var itemId interface{}
			if value == "" {
				itemId = nil
			} else {
				itemId = value
			}
			// not empty
			fmt.Println("key", key)
			fmt.Println("value", value)
			s.repo.UpdateItemToSetTx(tx, itemSetId, key, itemId)
		}

		return nil
	})
	// return &pb.UpdateItemSetsToBuildResponse{}, nil
}

/**
* Deletes a build for a member by its build id.
**/
func (s *service) DeleteBuildByMember(ctx context.Context, req *pb.DeleteBuildByMemberRequest) (*pb.DeleteBuildByMemberResponse, error) {
	// check if build is member's
	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}
	buildId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = s.GetBuildForMemberById(memberId, buildId)

	if err != nil {
		return nil, fmt.Errorf("The build does not belong to this member or does not exist.")
	}

	// delete build from db
	err = s.repo.DeleteBuildByIdForMember(memberId, buildId)

	if err != nil {
		return nil, err
	}

	return &pb.DeleteBuildByMemberResponse{}, nil
}
