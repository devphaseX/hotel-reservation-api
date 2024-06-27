package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type signInBodyParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResp struct {
	User  types.User `json:"user"`
	Token string     `json:"token"`
}

type FailedResp struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(FailedResp{
		Type:    "error",
		Message: "invalid credentials mismatch email or password",
	})
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var body signInBodyParams

	if err := c.BodyParser(&body); err != nil {
		fmt.Println(err)
		return errors.New("invalid payload received")
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), body.Email)

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}

		fmt.Println(err)
		return errors.New("failed to sign user in")
	}

	if !types.IsValidPassword(user.EncryptedPassword, body.Password) {
		return invalidCredentials(c)
	}

	token, err := CreateTokenClaim(user)
	if err != nil {
		fmt.Println("failed to sign token: %w", err)
		return errors.New("failed to sign user")
	}

	return c.JSON(SignInResp{User: *user, Token: token})
}

func CreateTokenClaim(user *types.User) (string, error) {
	expires := time.Now().Add(time.Hour * 4)
	claim := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	secret := os.Getenv("JWT_SECRET")
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
}
