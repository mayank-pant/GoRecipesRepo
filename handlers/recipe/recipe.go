package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	cache "recipes/cache"
	"recipes/datastore"
	"recipes/handlers"
	authenticate "recipes/models/authentication"
	recipes "recipes/models/recipe"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func SearchRecipes(w http.ResponseWriter, r *http.Request) {
	nameParam := r.FormValue("name")
	prepTimeParam := r.FormValue("preptime")
	difficultyParam := r.FormValue("difficulty")
	vegetarianParam := r.FormValue("vegetarian")

	var params = make(map[string]interface{})

	if len(difficultyParam) > 0 {
		difficulty, err := strconv.Atoi(difficultyParam)
		if err != nil {
			handlers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		params["difficulty"] = difficulty
	}

	if len(vegetarianParam) > 0 {
		vegetarian, err := strconv.ParseBool(vegetarianParam)
		if err != nil {
			handlers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		params["vegetarian"] = vegetarian
	}

	if len(prepTimeParam) > 0 {
		prepTime, err := strconv.Atoi(prepTimeParam)
		if err != nil {
			handlers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		params["preptime"] = prepTime
	}

	if len(nameParam) > 0 {
		params["name"] = nameParam
	}

	dataStoreClient, err := datastore.GetDatastoreClient(datastore.POSTGRE)
	if err != nil {
		log.Println("Error in obtaining datastore client")
		log.Println(err.Error())
		return
	}

	defer dataStoreClient.Close()

	searchResults, err := dataStoreClient.SearchRecipes(params, r.URL.Query())
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, searchResults)

}

func RateRecipe(w http.ResponseWriter, r *http.Request) {
	dataStoreClient, err := datastore.GetDatastoreClient(datastore.POSTGRE)
	if err != nil {
		log.Println("Error in obtaining datastore client")
		log.Println(err.Error())
		handlers.RespondWithError(w, http.StatusInternalServerError, "OOPS! Something went wrong. Please try again a bit later")
		return
	}

	defer dataStoreClient.Close()

	var ratingGiven = &recipes.Rating{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ratingGiven); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	vars := mux.Vars(r)
	recipeID, _ := strconv.Atoi(vars["id"])

	_, err = dataStoreClient.RateRecipe(recipeID, ratingGiven.Rating)
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update in cache
	cacheClient, err := cache.GetCacheClient(cache.REDIS)
	if err != nil {
		log.Println(err)
	}
	defer cacheClient.Close()

	isDeleted, err := cacheClient.Delete(string(recipeID))
	if err != nil || !isDeleted {
		log.Println("Failed to delete in cache")
	}

	handlers.RespondWithJSON(w, http.StatusOK, "success")

}

func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe = &recipes.Recipe{}

	TokenClaim := context.Get(r, "decoded")

	decodedTokenClaims, ok := TokenClaim.(authenticate.TokenClaims)
	if !ok {
		log.Println("Failed to validate user id")
		handlers.RespondWithError(w, http.StatusInternalServerError, "Failed to validate user id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recipe); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	// Set userId validated through the token istead of from the user
	recipe.UserId = decodedTokenClaims.UserId

	dataStoreClient, err := datastore.GetDatastoreClient(datastore.POSTGRE)
	if err != nil {
		log.Println("Error in obtaining datastore client")
		log.Println(err.Error())
		handlers.RespondWithError(w, http.StatusInternalServerError, "OOPS! Something went wrong. Please try again a bit later")
		return
	}

	defer dataStoreClient.Close()

	recipeId, err := dataStoreClient.CreateRecipe(*recipe)

	if err != nil {
		log.Println(err.Error())
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	cacheClient, err := cache.GetCacheClient(cache.REDIS)
	if err != nil {
		log.Println(err)
	}
	defer cacheClient.Close()
	// Add recipe to cache for faster fetching
	b, err := json.Marshal(recipe)
	if err != nil {
		log.Println("Error marshalling recipe")
		log.Println(err.Error())
	} else {
		err = cacheClient.Set(string(recipeId), b)
		if err != nil {
			log.Println("Error caching")
			log.Println(err.Error())
		}
	}

	handlers.RespondWithJSON(w, http.StatusOK, map[string]int{"recipeId": recipeId})

}

// GetRecipe used to get recipe detail
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Error in obtaining recipeId")
		handlers.RespondWithError(w, http.StatusInternalServerError, "Invalid recipeId")
		return
	}

	cacheClient, err := cache.GetCacheClient(cache.REDIS)
	if err != nil {
		log.Println(err)
	}

	if cacheClient != nil {
		defer cacheClient.Close()
		var cachedRecipe = &recipes.Recipe{}
		cachedBytes, err := cacheClient.Get(string(recipeID))
		if err != nil {
			log.Println("Could not get recipe from cache")
		} else {
			err = json.Unmarshal([]byte(cachedBytes), cachedRecipe)
			if err != nil {
				log.Println(err.Error())
			} else {
				cachedRecipe.UId = recipeID
				handlers.RespondWithJSON(w, http.StatusOK, cachedRecipe)
				return
			}
		}
	}

	dataStoreClient, err := datastore.GetDatastoreClient(datastore.POSTGRE)
	if err != nil {
		log.Println("Error in obtaining datastore client")
		log.Println(err.Error())
		handlers.RespondWithError(w, http.StatusInternalServerError, "OOPS! Something went wrong. Please try again a bit later")
		return
	}

	defer dataStoreClient.Close()

	recipe, err := dataStoreClient.GetRecipe(recipeID)
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update in cache the recipe data. Fetching client as it was deferred before
	cacheClient, err = cache.GetCacheClient(cache.REDIS)
	if err != nil {
		log.Println(err)
	}
	defer cacheClient.Close()
	// Add recipe to cache for faster fetching
	b, err := json.Marshal(recipe)
	if err != nil {
		log.Println("Error marshalling recipe")
		log.Println(err.Error())
	} else {
		err = cacheClient.Set(string(recipeID), b)
		if err != nil {
			log.Println("Error caching")
			log.Println(err.Error())
		}
	}

	handlers.RespondWithJSON(w, http.StatusOK, recipe)

}

// UpdateRecipe to update an exisitng recipe
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	dataStoreClient, err := datastore.GetDatastoreClient(datastore.POSTGRE)
	if err != nil {
		log.Println("Error in obtaining datastore client")
		log.Println(err.Error())
		handlers.RespondWithError(w, http.StatusInternalServerError, "OOPS! Something went wrong. Please try again a bit later")
		return
	}

	defer dataStoreClient.Close()

	vars := mux.Vars(r)
	recipeID, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Println("Error in obtaining recipeId")
		handlers.RespondWithError(w, http.StatusInternalServerError, "Invalid recipeId")
		return
	}

	TokenClaim := context.Get(r, "decoded")

	decodedTokenClaims, ok := TokenClaim.(authenticate.TokenClaims)
	if !ok {
		log.Println("Failed to validate user id")
		handlers.RespondWithError(w, http.StatusInternalServerError, "Failed to validate user id")
		return
	}

	var recipe = &recipes.Recipe{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recipe); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	// Set userId validated through the token istead of from the user
	recipe.UserId = decodedTokenClaims.UserId
	// Recipe Id is taken from the API path
	recipe.UId = recipeID

	_, err = dataStoreClient.UpdateRecipe(*recipe)
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Update in cache
	cacheClient, err := cache.GetCacheClient(cache.REDIS)
	if err != nil {
		log.Println(err)
	}
	defer cacheClient.Close()
	// Add recipe to cache for faster fetching
	b, err := json.Marshal(recipe)
	if err != nil {
		log.Println("Error marshalling recipe")
		log.Println(err.Error())
	} else {
		err = cacheClient.Set(string(recipeID), b)
		if err != nil {
			log.Println("Error caching")
			log.Println(err.Error())
		}
	}

	handlers.RespondWithJSON(w, http.StatusOK, "success")
}
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	dataStoreClient, err := datastore.GetDatastoreClient(datastore.POSTGRE)
	if err != nil {
		log.Println("Error in obtaining datastore client")
		log.Println(err.Error())
		handlers.RespondWithError(w, http.StatusInternalServerError, "OOPS! Something went wrong. Please try again a bit later")
		return
	}

	defer dataStoreClient.Close()

	vars := mux.Vars(r)
	recipeId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Error in obtaining recipeId")
		handlers.RespondWithError(w, http.StatusInternalServerError, "Invalid recipeId")
		return
	}

	TokenClaim := context.Get(r, "decoded")

	decodedTokenClaims, ok := TokenClaim.(authenticate.TokenClaims)
	if !ok {
		log.Println("Failed to validate user id")
		handlers.RespondWithError(w, http.StatusInternalServerError, "Failed to validate user id")
		return
	}

	existingRecipe, err := dataStoreClient.GetRecipe(recipeId)
	if err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, "Recipe doesn't exist")
		return
	}

	if existingRecipe.UserId != decodedTokenClaims.UserId {
		handlers.RespondWithError(w, http.StatusBadRequest, "User is not authorized to perform this action")
		return
	}

	_, err = dataStoreClient.DeleteRecipe(recipeId)

	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update in cache
	cacheClient, err := cache.GetCacheClient(cache.REDIS)
	if err != nil {
		log.Println(err)
	}
	defer cacheClient.Close()

	isDeleted, err := cacheClient.Delete(string(recipeId))
	if err != nil || !isDeleted {
		log.Println("Failed to delete in cache")
	}

	handlers.RespondWithJSON(w, http.StatusOK, "success")

}

// Listrecipes lists all recipes
func ListRecipes(w http.ResponseWriter, r *http.Request) {
	dataStoreClient, err := datastore.GetDatastoreClient(datastore.POSTGRE)
	if err != nil {
		log.Println("Error in obtaining datastore client")
		log.Println(err.Error())
		handlers.RespondWithError(w, http.StatusInternalServerError, "OOPS! Something went wrong. Please try again a bit later")
		return
	}

	defer dataStoreClient.Close()

	recipeList, err := dataStoreClient.ListRecipes(r.URL.Query())
	if err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	handlers.RespondWithJSON(w, http.StatusOK, recipeList)
}
