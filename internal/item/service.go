package item

import (
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
)

type ItemService struct {
	Repo *ItemRepository
}

func NewItemService(repo *ItemRepository) *ItemService {
	return &ItemService{
		Repo: repo,
	}
}

func (s *ItemService) CreateItemService(createItemReq CreateItemRequest) error {
	return s.Repo.CreateItem(createItemReq)
}

func (s *ItemService) AddItemToBuildService(memberId uuid.UUID, item CreateItemRequest) error {
	return s.Repo.CreateItem(item)
}

func (s *ItemService) GetItemsService() (*[]models.Item, error) {
	return s.Repo.GetItems()
}

func (s *ItemService) UpdateItemsService(id uuid.UUID, updateItemReq UpdateItemReq) (*models.Item, error) {
	return s.Repo.UpdateItemById(id, updateItemReq)
}

func (s *ItemService) GetWikiItemsService() (*[]WikiItem, error) {
	return s.Repo.GetWikiItems()
}

func (s *ItemService) GetBaseItemsService() (*[]BaseItem, error) {
	return s.Repo.GetBaseItems()
}
