package api

import (
	"errors"
	"fmt"

	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().Value("user").(*types.User)
	fmt.Println("user", user)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	return user, nil
}
