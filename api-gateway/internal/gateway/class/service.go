package class

import "fmt"

/**
* Classes an Ascendancies Services
**/
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
func (s *ClassService) CreateDefaultClassesAndAscendanciesService(classes []CreateDefaultClass, ascendancies []CreateDefaultAscendancy) error {

	if err := s.Repo.BatchCreateDefaultClasses(classes); err != nil {
		fmt.Println("Error when creating class:", err)
		return err
	}

	if err := s.Repo.BatchCreateDefaultAscendancies(ascendancies); err != nil {
		fmt.Println("Error when creating ascendancy:", err)
		return err
	}

	return nil
}

/**
* Gets list of classes and ascendancies.
**/
func (s *ClassService) GetClassesAndAscendancies() (*GetClassesAndAscendanciesResponse, error) {

	return s.Repo.GetClassesAndAscendancies()
}
