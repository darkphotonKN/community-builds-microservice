package item

import (
	"fmt"

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

func (s *ItemService) CreateItemService(memberId uuid.UUID, item models.Item) error {
	fmt.Println("Creating item with memberId:", memberId)
	// creating item for a specific user
	item.MemberID = memberId

	// generate a new product id
	newProdId := uuid.New()
	item.ProductID = newProdId
	return s.Repo.CreateItem(item)
}

func (s *ItemService) GetItemsService(memberId uuid.UUID) (*[]models.Item, error) {
	return s.Repo.GetItems(memberId)
}

func (s *ItemService) UpdateItemsService(memberId uuid.UUID, id uuid.UUID, updateItemReq UpdateItemReq) (*models.Item, error) {
	return s.Repo.UpdateItemById(memberId, id, updateItemReq)
}
