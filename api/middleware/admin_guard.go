package middleware

import (
	"net/http"

	"github.com/devphaseX/hotel-reservation-api/api"
	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)

	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(api.FailedResp{Type: "error", Message: "unauthorized"})
	}

	if !user.IsAdmin {
		return c.Status(http.StatusForbidden).JSON(api.FailedResp{Type: "error", Message: "user not an admin"})
	}

	return c.Next()
}
