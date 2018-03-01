package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	auth "recipes/models/authentication"
	recipes "recipes/models/recipe"
	"testing"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestingSignUp(t *testing.T) {
	var user = &auth.Authentication{Username: "sandybhat", Password: "test"}
	payloadBytes, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payloadBytes))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestingAuthentication(t *testing.T) {
	var user = &auth.Authentication{Username: "sandybhat2", Password: "test"}
	payloadBytes, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payloadBytes))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("POST", "/authenticate", bytes.NewBuffer(payloadBytes))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestingCreateRecipe(t *testing.T) {
	var user = &auth.Authentication{Username: "sandybhat2", Password: "test"}
	payloadBytes, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewBuffer(payloadBytes))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	token := response.Body.String()
	var newRecipe = &recipes.Recipe{Name: "testrecipe", PrepTime: 12, Difficulty: 3, Vegetarian: true}
	recipePayloadBytes, _ := json.Marshal(newRecipe)
	req, _ = http.NewRequest("POST", "/recipes", bytes.NewBuffer(recipePayloadBytes))
	req.Header.Add("Authorization", "Buffer "+token)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestingListRecipe(t *testing.T) {
	req, _ := http.NewRequest("GET", "/recipes", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}
