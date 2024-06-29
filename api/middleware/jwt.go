package middleware

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JWTAuth(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]

		if !ok {
			return utils.ErrUnauthorized("api token not present")
		}

		claim, err := ParseJWT(token[0])
		if err != nil {
			return err
		}

		expires, err :=
			time.Parse(time.RFC3339, claim["expires"].(string))

		if err != nil || expires.Before(time.Now()) {
			return errors.New("expired token")
		}

		id, err := primitive.ObjectIDFromHex(claim["id"].(string))

		if err != nil {
			return err
		}

		if ok {
			user, err := userStore.GetUserById(c.Context(), id)

			if err != nil {
				return err
			}

			c.Context().SetUserValue("user", user)
		}

		return c.Next()
	}
}

func ParseJWT(tokenStr string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method:", token.Header["alg"])
			return nil, utils.ErrUnauthorized()
		}

		secret := os.Getenv("JWT_SECRET")

		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse invalid token:", err)
		return nil, utils.ErrUnauthorized()
	}

	if !token.Valid {
		return nil, utils.ErrUnauthorized("token not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
