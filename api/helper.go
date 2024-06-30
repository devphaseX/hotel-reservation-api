package api

import (
	"errors"

	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	return user, nil
}

type ResourceResp struct {
	Data    any   `json:"data"`
	Results int64 `json:"results"`
	Page    int64 `json:"page"`
}

func NewResourceResp(data any, results, page int64) ResourceResp {
	return ResourceResp{
		Data:    data,
		Results: results,
		Page:    page,
	}
}

// type Filter map[string]any

// func NewFilter(q map[string]any, filterKeys []string) Filter {
// 	if len(filterKeys) == 0 {
// 		return Filter{}
// 	}

// 	cq := utils.CloneMap(q)
// 	for _, key := range filterKeys {
// 		delete(cq, key)
// 	}

// 	return Filter(cq)
// }
