package main

import (
	"bytes"
	ctx "context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"backend/api"
	"backend/models"

	assert "github.com/stretchr/testify/assert"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/stevens-tyr/tyr-gin"
)

func cleanUserDB(email string) {
	db, _ := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	assignCol := tyrgin.GetMongoCollection("users", db)

	assignCol.DeleteOne(ctx.Background(), bson.M{"email": email})
}

func performRequest(r http.Handler, method, path string, body []byte, token string) *httptest.ResponseRecorder {
	var req *http.Request

	if method == "GET" {
		req, _ = http.NewRequest(method, path, nil)
	} else {
		bb := new(bytes.Buffer)
		json.NewEncoder(bb).Encode(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	if token != "" {
		bearer := "Bearer " + token
		req.Header.Add("Authorization", bearer)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}

func TestInvalidLogin(t *testing.T) {
	server := api.SetUp()

	cleanUserDB("test@test.com")
	loginData := models.Login{
		Email:    "test@test.com",
		Password: "tester123",
	}
	jsb, _ := json.Marshal(loginData)

	resp := performRequest(server, "POST", "/api/v1/auth/login", jsb, "")
	assert.Equal(t, 401, resp.Code)
}

func TestInvalidRegister(t *testing.T) {
	server := api.SetUp()
	cleanUserDB("test@test.com")

	registerData := models.RegisterForm{
		Email:                "test@test.com",
		Password:             "tester123",
		PasswordConfirmation: "tester1234",
		First:                "Bob",
		Last:                 "Bobert",
	}
	jsb, _ := json.Marshal(registerData)

	resp := performRequest(server, "POST", "/api/v1/auth/register", jsb, "")
	assert.Equal(t, 400, resp.Code)
}

func TestTakenRegister(t *testing.T) {
	server := api.SetUp()
	cleanUserDB("test2@test.com")

	registerData := models.RegisterForm{
		Email:                "test2@test.com",
		Password:             "tester123",
		PasswordConfirmation: "tester123",
		First:                "Bob",
		Last:                 "Bobert",
	}
	jsb, _ := json.Marshal(registerData)

	resp := performRequest(server, "POST", "/api/v1/auth/register", jsb, "")
	assert.Equal(t, 200, resp.Code)

	registerData = models.RegisterForm{
		Email:                "test2@test.com",
		Password:             "tester123",
		PasswordConfirmation: "tester123",
		First:                "Bob",
		Last:                 "Bobert",
	}
	jsb, _ = json.Marshal(registerData)

	resp = performRequest(server, "POST", "/api/v1/auth/register", jsb, "")
	assert.Equal(t, 400, resp.Code)

	cleanUserDB("test2@test.com")
}

func TestValidRegisterAndLogin(t *testing.T) {
	server := api.SetUp()
	cleanUserDB("test@test.com")

	registerData := models.RegisterForm{
		Email:                "test@test.com",
		Password:             "tester123",
		PasswordConfirmation: "tester123",
		First:                "Bob",
		Last:                 "Bobert",
	}
	jsb, _ := json.Marshal(registerData)

	resp := performRequest(server, "POST", "/api/v1/auth/register", jsb, "")
	assert.Equal(t, 200, resp.Code)

	loginData := models.Login{
		Email:    "test@test.com",
		Password: "tester123",
	}
	jsb, _ = json.Marshal(loginData)

	resp = performRequest(server, "POST", "/api/v1/auth/login", jsb, "")
	assert.Equal(t, 200, resp.Code)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)

	resp = performRequest(server, "GET", "/api/v1/auth/logged_in", nil, body["token"].(string))
	assert.Equal(t, 200, resp.Code)

	cleanUserDB("test@test.com")
}
