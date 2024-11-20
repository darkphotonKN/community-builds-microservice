package tag

import (
	"github.com/darkphotonKN/community-builds/internal/models"
)

type TagService struct {
	Repo *TagRepository
}

func NewTagService(repo *TagRepository) *TagService {
	return &TagService{
		Repo: repo,
	}
}

func (s *TagService) CreateTagService(createTagReq CreateTagRequest) error {
	return s.Repo.CreateTag(createTagReq)
}

func (s *TagService) UpdateTagsService(updateTagReq UpdateTagRequest) error {
	return s.Repo.UpdateTag(updateTagReq)
}

func (s *TagService) GetTagsService() (*[]models.Tag, error) {
	return s.Repo.GetTags()
}
