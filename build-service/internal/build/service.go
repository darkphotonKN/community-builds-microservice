package build

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/skill"
	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/tag"
	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/build"
	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/types"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	repo         Repository
	publishCh    *amqp.Channel
	skillService skill.Service
	tagService   tag.Service
}

type Repository interface {
	CreateBuild(memberId uuid.UUID, createBuildRequest CreateBuildRequest) (*uuid.UUID, error)
	GetBuildsByMemberId(id uuid.UUID) (*[]BuildListQuery, error)
	GetAllBuilds(pageNo int, pageSize int, sortOrder string, sortBy string, search string, skillId uuid.UUID, minRating *int, ratingCategory types.RatingCategory) (*[]BuildListQuery, error)
	GetBuildInfo(buildId uuid.UUID) (*BuildInfoResponse, error)
	GetAndFormSkillLinks(skillData []models.SkillRow) SkillGroupResponse
}

func NewService(repo Repository, publishCh *amqp.Channel, skillService skill.Service, tagService tag.Service) Service {
	return &service{repo: repo, publishCh: publishCh, skillService: skillService, tagService: tagService}
}

/**
* Get list of builds available to a member.
**/
func (s *service) GetBuildsForMember(memberId uuid.UUID) (*[]BuildListResponse, error) {
	baseBuilds, err := s.repo.GetBuildsByMemberId(memberId)

	if err != nil {
		return nil, err
	}

	buildListResponse := make([]BuildListResponse, len(*baseBuilds))

	// query and add each builds tag's
	for index, build := range *baseBuilds {
		tags, err := s.tagService.GetBuildTagsForMemberById(build.Id)

		// stop query pre-maturely if errored on query
		if err != nil {
			return nil, err
		}

		buildListResponse[index] = BuildListResponse{
			Id:                 build.Id,
			Title:              build.Title,
			Description:        build.Description,
			Class:              build.Class,
			Ascendancy:         build.Ascendancy,
			MainSkillName:      build.MainSkillName,
			AvgEndGameRating:   build.AvgEndGameRating,
			AvgBossingRating:   build.AvgBossingRating,
			AvgCreativeRating:  build.AvgCreativeRating,
			AvgFunRating:       build.AvgFunRating,
			AvgSpeedFarmRating: build.AvgSpeedFarmRating,
			Tags:               *tags,
			Views:              build.Views,
			Status:             build.Status,
			CreatedAt:          build.CreatedAt,
		}
	}

	return &buildListResponse, nil
}

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
		tags, err := s.tagService.GetBuildTagsForMemberById(build.Id)

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
	builds, err := s.GetBuildsForMember(memberId)

	if err != nil {
		return nil, err
	}

	if len(*builds) > maxBuildCount {
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
	// err = s.repo.CreateBuildTags(*buildId, createBuildRequest.TagIDs)

	// if err != nil {
	// 	return nil, err
	// }

	// create build default set
	// err = s.CreateDefaultItemSetsToBuildService(memberId, *buildId)

	// if err != nil {
	// 	return nil, err
	// }

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

func (s *service) GetBuildsByMemberId(ctx context.Context, req *pb.GetBuildsByMemberIdRequest) (*pb.GetBuildsByMemberIdResponse, error) {
	memberId, err := uuid.Parse(req.MemberId)

	if err != nil {
		return nil, err
	}

	baseBuilds, err := s.repo.GetBuildsByMemberId(memberId)

	if err != nil {
		return nil, err
	}

	pbBuilds := make([]*pb.BuildList, len(*baseBuilds))

	// query and add each builds tag's
	for index, build := range *baseBuilds {
		tags, err := s.tagService.GetBuildTagsForMemberById(build.Id)

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

	return &pb.GetBuildsByMemberIdResponse{
		Builds: pbBuilds,
	}, nil
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
		Id:          info.ID.String(),
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
