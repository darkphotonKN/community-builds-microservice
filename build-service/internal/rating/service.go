package rating

import (
	"context"
	"fmt"

	// "github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/build"
	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/build"
	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/rating"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/types"
	"github.com/google/uuid"
)

type service struct {
	repo         Repository
	buildService build.Service
}

type Repository interface {
	CreateRatingForBuildById(request CreateRating) error
	GetAllRatingsByMemberId(memberId uuid.UUID) (*[]models.Rating, error)
	GetAllRatingsByCategoryForBuild(buildId string, category types.RatingCategory) ([]int, error)
}

func NewService(repo Repository, buildService build.Service) Service {
	return &service{repo: repo, buildService: buildService}
}

/**
* Posts single rating for a single product.
**/
func (s *service) CreateRatingByBuildId(ctx context.Context, req *pb.CreateRatingByBuildIdRequest) (*pb.CreateRatingByBuildIdResponse, error) {
	// create rating for build
	fmt.Println("CreateRatingByBuildId start", req)
	fromMemberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}

	// get build by id to find the owner (toMemberId)
	buildId, err := uuid.Parse(req.BuildId)
	if err != nil {
		return nil, err
	}

	build, err := s.buildService.GetBuildById(buildId)

	request := CreateRating{
		FromMemberId: fromMemberId,
		ToMemberId:   build.MemberID,
		BuildId:      build.ID,
		Rating:       int(req.Value),
	}
	fmt.Printf("Creating rating for build %+v through service\n", request)
	err = s.repo.CreateRatingForBuildById(request)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

/**
* TODO: need to change from products to builds.
* Gets all Ratings.
**/
func (s *service) GetAllRatingsForProduct(userId uuid.UUID) (*[]models.Rating, error) {
	return s.repo.GetAllRatingsByMemberId(userId)
}
