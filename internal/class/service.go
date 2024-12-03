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
func (s *ClassService) CreateClassesAndAscendanciesService(classes []CreateClass) error {

	for _, class := range classes {
		err := s.Repo.CreateClass(class)
		fmt.Println("Error when creating class:", err)
		if err != nil {
			return err
		}
	}

	fmt.Println("Successfully created all classes!")
	return nil
}
