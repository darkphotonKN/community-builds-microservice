package build

type BuildHandler struct {
	Service *BuildService
}

func NewBuildHandler(service *BuildService) *BuildHandler {
	return &BuildHandler{
		Service: service,
	}
}
