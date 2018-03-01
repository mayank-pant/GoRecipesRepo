package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"recipes/authentication"
	"recipes/datastore"
	"recipes/handlers"
	authenticate "recipes/models/authentication"
)

// SignUp for a new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	var auth authenticate.Authentication
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&auth); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	dataStoreClient, err := datastore.GetDatastoreClient(datastore.POSTGRE)
	if err != nil {
		log.Println("Error in obtaining datastore client")
		log.Println(err.Error())
		handlers.RespondWithError(w, http.StatusInternalServerError, "OOPS! Something went wrong. Please try again a bit later")
		return
	}

	defer dataStoreClient.Close()

	user, err := dataStoreClient.FindUser(auth.Username)
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if user != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := dataStoreClient.CreateUser(auth.Username, auth.Password)

	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !created {
		handlers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, "SUCCESS")

}

// ObtainToken to get authorization token
func ObtainToken(w http.ResponseWriter, r *http.Request) {

	var auth authenticate.Authentication
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&auth); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	dataStoreClient, err := datastore.GetDatastoreClient(datastore.POSTGRE)
	if err != nil {
		log.Println("Error in obtaining datastore client")
		log.Println(err.Error())
		handlers.RespondWithError(w, http.StatusInternalServerError, "OOPS! Something went wrong. Please try again a bit later")
		return
	}

	defer dataStoreClient.Close()

	userId, err := dataStoreClient.ValidateUser(auth.Username, auth.Password)

	if err != nil || userId == 0 {
		handlers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token := authentication.CreateToken(auth.Username, userId)
	handlers.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
