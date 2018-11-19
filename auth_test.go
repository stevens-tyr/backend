package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/api"
	"backend/models"

	assert "github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
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
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}

func TestInvalidLogin(t *testing.T) {
	server := api.SetUp()
	server.Run(":5555")

	loginData := models.Login{
		Email:    "test@test.com",
		Password: "tester123",
	}
	jsb, _ := json.Marshal(loginData)

	resp := performRequest(server, "POST", "/api/v1/auth/login", jsb)
	assert.Equal(t, 401, resp.Code)
}
