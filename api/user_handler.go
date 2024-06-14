package api

import (
	"errors"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserById(c.Context(), id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "user not found"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandlerGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())

	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) != 0 {
		return c.JSON(map[string]any{"error": errors})
	}

	user, err := types.NewUserFromParams(params)

	if err != nil {
		return err
	}

	newUser, err := h.userStore.CreateUser(c.Context(), user)

	if err != nil {
		return err
	}

	return c.JSON(newUser)
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		userId = c.Params("id")
		values = bson.M{}
	)

	oid, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return err
	}

	if err = c.BodyParser(&values); err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	if err = h.userStore.PutUser(c.Context(), filter, values); err != nil {
		return err
	}

	return c.JSON(map[string]string{"message": "user recorded updated"})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	err := h.userStore.RemoveUser(c.Context(), userId)

	if err != nil {
		return err
	}

	return c.JSON(map[string]string{"message": "user deleted", "deleted": userId})
}
