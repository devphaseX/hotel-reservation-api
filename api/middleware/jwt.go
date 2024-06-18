package middleware

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]

	if !ok {
		return errors.New("authorized")
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

	return c.Next()
}

func ParseJWT(tokenStr string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method:", token.Header["alg"])
			return nil, errors.New("Unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")

		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse invalid token:", err)
		return nil, errors.New("Unauthorized")
	}

	if !token.Valid {
		return nil, errors.New("token not valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
