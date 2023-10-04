package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/r4chi7/aspire-lite/controller"
	"github.com/r4chi7/aspire-lite/database"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/r4chi7/aspire-lite/repository"
	"github.com/r4chi7/aspire-lite/service"
)

func TestUser(t *testing.T) {
	db := database.Get()
	userRepository := repository.NewUser(db)
	userController := controller.NewUser(service.NewUser(userRepository))
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"email":"test@example.com","password":"test@123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userController.Create)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := resp["token"]; !exists {
		t.Error("expected response body to have token field")
	}

	var insertedUser model.User
	db.Where("email = ?", "test@example.com").Find(&insertedUser)
	if insertedUser.ID == 0 {
		t.Error("expected user to be created")
	}
}
