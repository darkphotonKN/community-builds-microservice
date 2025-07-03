package class

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/class"
)

/**
* Classes an Ascendancies Services
**/

type service struct {
	repo Repository
}

type Repository interface {
	BatchCreateDefaultClasses(classes []CreateDefaultClass) error
	BatchCreateDefaultAscendancies(ascendancies []CreateDefaultAscendancy) error
	GetClassesAndAscendancies() (*GetClassesAndAscendanciesResponse, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

/**
* Creates all base classes and ascendancies
**/
// func (s *service) CreateDefaultClassesAndAscendancies(classes []CreateDefaultClass, ascendancies []CreateDefaultAscendancy) error {

// 	if err := s.repo.BatchCreateDefaultClasses(classes); err != nil {
// 		fmt.Println("Error when creating class:", err)
// 		return err
// 	}

// 	if err := s.repo.BatchCreateDefaultAscendancies(ascendancies); err != nil {
// 		fmt.Println("Error when creating ascendancy:", err)
// 		return err
// 	}

// 	return nil
// }

/**
* Gets list of classes and ascendancies.
**/
func (s *service) GetClassesAndAscendancies(ctx context.Context, req *pb.GetClassesAndAscendanciesRequest) (*pb.GetClassesAndAscendanciesResponse, error) {
	res, err := s.repo.GetClassesAndAscendancies()
	if err != nil {
		return nil, err
	}
	pbClasses := make([]*pb.Class, len(res.Classes))
	for _, class := range res.Classes {
		pbClass := &pb.Class{
			Id:          class.ID.String(),
			Name:        class.Name,
			Description: class.Description,
			ImageUrl:    class.ImageURL,
			CreatedAt:   class.CreatedAt.String(),
			UpdatedAt:   class.UpdatedAt.String(),
		}
		pbClasses = append(pbClasses, pbClass)
	}

	pbAscendancies := make([]*pb.Ascendancy, len(res.Ascendancies))
	for _, Ascendancy := range res.Ascendancies {
		pbClass := &pb.Ascendancy{
			Id:          Ascendancy.ID.String(),
			ClassId:     Ascendancy.ClassID.String(),
			Name:        Ascendancy.Name,
			Description: Ascendancy.Description,
			ImageUrl:    Ascendancy.ImageURL,
			CreatedAt:   Ascendancy.CreatedAt.String(),
			UpdatedAt:   Ascendancy.UpdatedAt.String(),
		}
		pbAscendancies = append(pbAscendancies, pbClass)
	}

	grpcRes := &pb.GetClassesAndAscendanciesResponse{
		Classes:      pbClasses,
		Ascendancies: pbAscendancies,
	}
	return grpcRes, nil
}
