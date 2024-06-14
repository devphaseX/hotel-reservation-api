package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDbUri  = "mongodb://localhost:27017"
	testDbName = "hotel-reservation-test"
)

type testUserStore struct {
	db.UserStore
}

func (ts *testUserStore) tearDown(t *testing.T) {
	if err := ts.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup() *testUserStore {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDbUri))
	if err != nil {
		log.Fatal(err)
	}

	return &testUserStore{
		UserStore: db.NewMongoUserStore(client, testDbName),
	}
}

func TestCreateUser(t *testing.T) {
	db := setup()

	defer db.tearDown(t)

	app := fiber.New()

	userHandler := NewUserHandler(db.UserStore)

	app.Post("/", userHandler.HandleCreateUser)

	params := types.CreateUserParams{
		FirstName: "James",
		LastName:  "Foo",
		Email:     "foo@gmail.com",
		Password:  "allowtempPas12",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))

	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Error(err)
	}

	var user types.User

	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Error("expected the user id to be set.")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Error("expected the encrypted password to be included.")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName %s but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected lastName %s but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}
