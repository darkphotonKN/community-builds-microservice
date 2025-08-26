package rating

import (
	"context"
	"encoding/json"
	"fmt"

	// "github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/build"
	"github.com/darkphotonKN/community-builds-microservice/build-service/internal/build"
	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/rating"
	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/types"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	repo         Repository
	publishCh    *amqp.Channel
	buildService build.Service
}

type Repository interface {
	CreateRatingForBuildById(request CreateRating) error
	GetAllRatingsByMemberId(memberId uuid.UUID) (*[]models.Rating, error)
	GetAllRatingsByCategoryForBuild(buildId string, category types.RatingCategory) ([]int, error)
}

func NewService(repo Repository, publishCh *amqp.Channel, buildService build.Service) Service {
	return &service{repo: repo, publishCh: publishCh, buildService: buildService}
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

	// publish to message broker
	payload := commonconstants.RatingCreatedEventPayload{
		UserID: build.MemberID.String(),
	}

	marshalledPayload, err := json.Marshal(payload)

	err = s.publishCh.PublishWithContext(
		ctx,
		commonconstants.RatingCreatedEvent,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshalledPayload,
			// persist message
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		fmt.Println("Failed to publish rating created event:", err)
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
