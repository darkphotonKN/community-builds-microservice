package skill

type SkillService struct {
	Repo *SkillRepository
}

func NewSkillService(repo *SkillRepository) *SkillService {
	return &SkillService{
		Repo: repo,
	}
}

func (s *SkillService) CreateSkillService(createSkillReq CreateSkillRequest) error {
	return s.Repo.CreateSkill(createSkillReq)
}
