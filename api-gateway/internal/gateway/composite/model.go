package composite

import (
	"context"

	itemPb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	skillPb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/skill"
)

// --- Request ---

type GameDataResponse struct {
	Items  []*itemPb.Item   `json:"items,omitempty"`
	Skills []*skillPb.Skill `json:"skills,omitempty"`
}

type ChWrapper struct {
	Type string
	Data interface{}
	Err  error
}

type CompositeClient interface {
	GetGameData(ctx context.Context) (*GameDataResponse, error)
}
