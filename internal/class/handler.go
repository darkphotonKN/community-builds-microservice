package class

type ClassHandler struct {
	Service *ClassService
}

func NewClassHandler(service *ClassService) *ClassHandler {
	return &ClassHandler{
		Service: service,
	}
}
