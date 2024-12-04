package class

import "fmt"

type ClassService struct {
	Repo *ClassRepository
}

func NewClassService(repo *ClassRepository) *ClassService {
	return &ClassService{
		Repo: repo,
	}
}

/**
* Creates all base classes and ascendancies
**/
func (s *ClassService) CreateDefaultClassesAndAscendanciesService(classes []CreateDefaultClass) error {

	if err := s.Repo.BatchCreateDefaultClasses(classes); err != nil {
		fmt.Println("Error when creating class:", err)
		return err
	}

	return nil
}
