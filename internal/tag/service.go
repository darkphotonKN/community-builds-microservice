package tag

import (
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
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

func (s *TagService) UpdateTagsService(id uuid.UUID, updateItemReq UpdateTagRequest) (*UpdateTagRequest, error) {
	return s.Repo.UpdateTag(id, updateItemReq)
}

func (s *TagService) GetTagsService() (*[]models.Tag, error) {
	return s.Repo.GetTags()
}
