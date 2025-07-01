package build

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/build"
	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/types"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	repo      Repository
	publishCh *amqp.Channel
}

type Repository interface {
	CreateBuild(memberId uuid.UUID, createBuildRequest CreateBuildRequest) (*uuid.UUID, error)
	GetBuildsByMemberId(id uuid.UUID) (*[]BuildListQuery, error)
	GetAllBuilds(pageNo int, pageSize int, sortOrder string, sortBy string, search string, skillId uuid.UUID, minRating *int, ratingCategory types.RatingCategory) ([]BuildListQuery, error)
}

func NewService(repo Repository, publishCh *amqp.Channel) Service {
	return &service{repo: repo, publishCh: publishCh}
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
		// tags, err := s.Repo.GetBuildTagsForMemberById(build.ID)

		// stop query pre-maturely if errored on query
		if err != nil {
			return nil, err
		}

		buildListResponse[index] = BuildListResponse{
			ID:                 build.Id,
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
			// Tags:               *tags,
			Views:     build.Views,
			Status:    build.Status,
			CreatedAt: build.CreatedAt,
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
func (s *service) GetCommunityBuildsService(pageNo int, pageSize int, sortOrder string, sortBy string, search string, skillId uuid.UUID, minRating *int, ratingCategory types.RatingCategory) ([]BuildListResponse, error) {

	builds, err := s.repo.GetAllBuilds(pageNo, pageSize, sortOrder, sortBy, search, skillId, minRating, ratingCategory)

	if err != nil {
		return nil, err
	}

	// get tags and add to build
	buildList := make([]BuildListResponse, len(builds))

	for index, build := range builds {
		// tags, err := s.repo.GetBuildTagsForMemberById(build.ID)

		// exit prematurely with error if any tags returned an error
		if err != nil {
			return nil, err
		}

		buildList[index] = BuildListResponse{
			ID:                 build.Id,
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
			// Tags:               *tags,
			Views:     build.Views,
			Status:    build.Status,
			CreatedAt: build.CreatedAt,
		}
	}

	return buildList, nil
}

func (s *service) CreateBuild(ctx context.Context, req *pb.CreateBuildRequest) (*pb.CreateBuildResponse, error) {
	memberId, err := uuid.Parse(req.MemberId)

	if err != nil {
		return nil, err
	}
	// confirm skill exists
	// _, err := s.SkillService.GetSkillByIdService(createBuildRequest.SkillID)

	// if err != nil {
	// 	return fmt.Errorf("main skill id could not be found when attempting to create build for it.")
	// }

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

	buildListResponse := make([]BuildListResponse, len(*baseBuilds))

	// query and add each builds tag's
	for index, build := range *baseBuilds {
		// tags, err := s.repo.GetBuildTagsForMemberById(build.ID)

		// stop query pre-maturely if errored on query
		if err != nil {
			return nil, err
		}

		buildListResponse[index] = BuildListResponse{
			ID:                 build.Id,
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
			// Tags:               *tags,
			Views:     build.Views,
			Status:    build.Status,
			CreatedAt: build.CreatedAt,
		}
	}

	// return &buildListResponse, nil
	return &pb.GetBuildsByMemberIdResponse{
		// Builds: pbBuilds,
	}, nil
}
