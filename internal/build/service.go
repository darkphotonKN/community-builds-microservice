package build

type BuildService struct {
	Repo *BuildRepository
}

func NewBuildService(repo *BuildRepository) *BuildService {
	return &BuildService{
		Repo: repo,
	}
}
