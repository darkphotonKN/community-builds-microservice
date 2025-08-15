package composite

import (
	"context"
	"fmt"

	itemPb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	skillPb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/skill"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	buildServiceName = "build-service"
	itemServiceName  = "item-service"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) CompositeClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) GetGameData(ctx context.Context) (*GameDataResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	buildServiceConn, err := discovery.ServiceConnection(ctx, buildServiceName, c.registry)
	itemServiceConn, err := discovery.ServiceConnection(ctx, itemServiceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer buildServiceConn.Close()
	defer itemServiceConn.Close()

	skillClient := skillPb.NewSkillServiceClient(buildServiceConn)

	itemClient := itemPb.NewItemServiceClient(itemServiceConn)

	ch := make(chan ChWrapper)

	go func() {

		skillResult, SkillErr := skillClient.GetSkills(ctx, &skillPb.GetSkillsRequest{})
		fmt.Println("skills", skillResult)
		if SkillErr != nil {
			fmt.Printf("Error when attempting to get skills: %s\n", SkillErr.Error())
			ch <- ChWrapper{
				Type: "skills",
				Err:  SkillErr,
			}
			return
		}
		ch <- ChWrapper{
			Type: "skills",
			Data: skillResult.Skills,
		}
	}()

	go func() {
		itemResult, itemErr := itemClient.GetItems(ctx, &itemPb.GetItemsRequest{
			Slot: "",
		})
		if itemErr != nil {
			fmt.Printf("Error when attempting to get items: %s\n", itemErr.Error())
			ch <- ChWrapper{
				Type: "items",
				Err:  itemErr,
			}
			return
		}
		ch <- ChWrapper{
			Type: "items",
			Data: itemResult.Items,
		}
	}()

	var items []*itemPb.Item
	var skills []*skillPb.Skill

	for i := 0; i < 2; i++ {
		select {
		case wrapper := <-ch:
			if wrapper.Type == "items" {
				if wrapper.Err != nil {
					return nil, fmt.Errorf("failed to get items: %w", wrapper.Err)
				}
				items = wrapper.Data.([]*itemPb.Item)
			}
			if wrapper.Type == "skills" {
				if wrapper.Err != nil {
					return nil, fmt.Errorf("failed to get skills: %w", wrapper.Err)
				}
				skills = wrapper.Data.([]*skillPb.Skill)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get game data: %w", err)
	}

	response := &GameDataResponse{
		Items:  items,
		Skills: skills,
	}
	fmt.Printf("Retrieved game data %+v through gateway after service discovery\n", response)

	return response, nil
}
