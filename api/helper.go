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
	Data    any `json:"data"`
	Results int `json:"results"`
	Page    int `json:"page"`
}

func NewResourceResp(data any, results, page int) ResourceResp {
	return ResourceResp{
		Data:    data,
		Results: results,
		Page:    page,
	}
}
