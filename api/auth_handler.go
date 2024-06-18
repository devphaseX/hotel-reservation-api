package api

import (
	"errors"
	"fmt"
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
	User  *types.User `json:"user"`
	Token string      `json:"token"`
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
			return errors.New("invalid credentials mismatch email or password")
		}

		fmt.Println(err)
		return errors.New("failed to sign user in")
	}

	if !types.IsValidPassword(user.EncryptedPassword, body.Password) {
		return errors.New("invalid credentials mismatch email or password")
	}

	token, err := createTokenClaim(user)

	if err != nil {
		fmt.Println("failed to sign token: %w", err)
		return errors.New("failed to sign user")
	}

	fmt.Println(token)
	return c.JSON(SignInResp{User: user, Token: token})
}

func createTokenClaim(user *types.User) (string, error) {
	expires := time.Now().Add(time.Hour * 4)
	claim := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	secret := os.Getenv("JWT_SECRET")
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
}
