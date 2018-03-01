package datastore

import (
	"net/url"
	recipes "recipes/models/recipe"
	users "recipes/models/user"
)

const (
	POSTGRE string = "postgre"
)

type DatastoreClient interface {
	Initialize() error
	Close() error
	FindUser(username string) (*users.User, error)
	CreateUser(username string, password string) (bool, error)
	ValidateUser(username string, password string) (int, error)
	CreateRecipe(userRecipe recipes.Recipe) (int, error)
	GetRecipe(recipeId int) (*recipes.Recipe, error)
	ListRecipes(urlValues url.Values) ([]recipes.Recipe, error)
	UpdateRecipe(userRecipe recipes.Recipe) (int, error)
	DeleteRecipe(recipeId int) (int, error)
	RateRecipe(recipeId int, rating int) (int, error)
	SearchRecipes(searchParams map[string]interface{}, urlValues url.Values) ([]recipes.Recipe, error)
}
