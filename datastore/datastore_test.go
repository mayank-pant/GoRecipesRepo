package datastore

import (
	recipe "recipes/models/recipe"
	"testing"
)

func TestCreateRecipe(t *testing.T) {
	datastoreClient, err := GetDatastoreClient(POSTGRE)
	if err != nil {
		t.Errorf(err.Error())
	}

	var newRecipe = recipe.Recipe{UId: 5555, UserId: 30, Name: "testrecipe", PrepTime: 30, Difficulty: 3, Vegetarian: true}
	_, err = datastoreClient.CreateRecipe(newRecipe)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetRecipe(t *testing.T) {
	datastoreClient, err := GetDatastoreClient(POSTGRE)
	if err != nil {
		t.Errorf(err.Error())
	}
	// Check for previously created recipe
	var recipeID = 555
	_, err = datastoreClient.GetRecipe(recipeID)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUpdateRecipe(t *testing.T) {
	datastoreClient, err := GetDatastoreClient(POSTGRE)
	if err != nil {
		t.Errorf(err.Error())
	}

	var newRecipe = recipe.Recipe{UId: 5555, UserId: 30, Name: "testrecipeupdated", PrepTime: 30, Difficulty: 3, Vegetarian: true}
	_, err = datastoreClient.UpdateRecipe(newRecipe)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestRateRecipe(t *testing.T) {
	datastoreClient, err := GetDatastoreClient(POSTGRE)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Rate previously created recipe
	var recipeID = 555
	_, err = datastoreClient.RateRecipe(recipeID, 4)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = datastoreClient.RateRecipe(recipeID, 6)
	// Invalid rating should throw error
	if err == nil {
		t.Errorf(err.Error())
	}
}
