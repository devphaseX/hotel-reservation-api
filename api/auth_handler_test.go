package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devphaseX/hotel-reservation-api/db/fixtures"
	"github.com/devphaseX/hotel-reservation-api/utils"
	"github.com/gofiber/fiber/v2"
)

func TestAuthSuccess(t *testing.T) {
	db := setup()

	defer db.tearDown(t)
	fixtures.AddUser(db.Store, "william", "stone", false)

	app := fiber.New()

	userHandler := NewAuthHandler(db.User)

	app.Post("/", userHandler.SignIn)

	params := signInBodyParams{
		Email:    "william.stone@mail.com",
		Password: "william_stone",
	}

	b, err := json.Marshal(params)

	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expect status code to be %d", http.StatusOK)
	}

	var authPayload SignInResp

	json.NewDecoder(resp.Body).Decode(&authPayload)

	if authPayload.Token == "" {
		t.Fatal("expected token to be present in auth response but got empty")
	}

	if authPayload.User.EncryptedPassword != "" {
		t.Fatal("expected encrypted password not to be pass to client")
	}

}

func TestAuthFailedWithWrongPassword(t *testing.T) {
	db := setup()

	defer db.tearDown(t)
	fixtures.AddUser(db.Store, "william", "stone", false)

	app := fiber.New()

	userHandler := NewAuthHandler(db.User)

	app.Post("/", userHandler.SignIn)

	params := signInBodyParams{
		Email:    "william.stone@gmail.com",
		Password: "wrongpassword",
	}

	b, err := json.Marshal(params)

	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("expect auth to failed to invalid credentials error")
	}

	var respPayload utils.FailedResp

	if err = json.NewDecoder(resp.Body).Decode(&respPayload); err != nil {
		t.Fatal(err)
	}

	if respPayload.Type != "error" {
		t.Fatalf("expect type of generated response to be of error type but got %v", respPayload.Type)
	}
}
