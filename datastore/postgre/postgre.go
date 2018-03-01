package datastore

import (
	"database/sql"
	"errors"
	"log"
	"net/url"
	"os"
	recipes "recipes/models/recipe"
	users "recipes/models/user"
	"strconv"
	"strings"

	"github.com/go-pg/pg/orm"

	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	//	DB_HOST     = "172.18.0.3"
	DB_USER     = "kapeiq"
	DB_PASSWORD = "kapeiq"
	DB_NAME     = "kapeiq"
)

type DatabaseClient struct {
	Handle *pg.DB
}

var postgreClient = &DatabaseClient{}

func GetClient() (*DatabaseClient, error) {
	err := postgreClient.Initialize()

	if err != nil {
		return nil, err
	}

	return postgreClient, nil
}

func (database *DatabaseClient) Initialize() error {
	DB_HOST := os.Getenv("DB_HOST")

	database.Handle = pg.Connect(&pg.Options{
		User:     "kapeiq",
		Password: "kapeiq",
		Database: "kapeiq",
		Addr:     DB_HOST + ":5432",
	})
	return nil
}

func (database *DatabaseClient) Close() error {
	err := database.Handle.Close()
	if err != nil {
		log.Println("Failed to close the handle")
		return err
	}
	return nil
}

func (database *DatabaseClient) FindUser(username string) (*users.User, error) {

	var user = &users.User{Username: username}

	// Check firs in Redis then DB. Use this method to validate the user before letting him alter something
	err := database.Handle.Model(user).Select()

	if err != nil && err != pg.ErrNoRows {
		log.Println(err.Error())
		return nil, err
	}

	if err == pg.ErrNoRows {
		log.Println("NO ROWS")
		return nil, nil
	}
	return user, nil

}

func (database *DatabaseClient) ValidateUser(username string, password string) (int, error) {
	user, err := database.FindUser(username)
	if err != nil && user != nil {
		log.Println(err.Error())
		return 0, err
	}

	if user == nil {
		return 0, nil
	}

	var userPassword = &users.Password{Id: user.Id}
	err = database.Handle.Select(userPassword)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	if err != nil {
		log.Println("Retrieving error")
		log.Println(err.Error())
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPassword.Hash), []byte(password))
	if err != nil {
		log.Println("hash compare error")
		log.Println(err.Error())
		return 0, err
	}

	return user.Id, nil
}

func (database *DatabaseClient) CreateUser(username string, password string) (bool, error) {

	var user = &users.User{Username: username, IsActive: true}

	err := database.Handle.Insert(user)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Println("Error generating hash")
		log.Println(err.Error())
		return false, err
	}

	var userPassword = &users.Password{Id: user.Id, Hash: string(hash)}

	err = database.Handle.Insert(userPassword)

	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	return true, nil

}

func (database *DatabaseClient) CreateRecipe(userRecipe recipes.Recipe) (int, error) {

	err := database.Handle.Insert(&userRecipe)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return userRecipe.UId, nil
}

func (database *DatabaseClient) GetRecipe(recipeId int) (*recipes.Recipe, error) {

	var recipe = &recipes.Recipe{UId: recipeId}

	err := database.Handle.Select(recipe)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return recipe, nil

}

func (database *DatabaseClient) ListRecipes(urlValues url.Values) ([]recipes.Recipe, error) {

	var recipeList = []recipes.Recipe{}

	if urlValues != nil {
		err := database.Handle.Model(&recipeList).Apply(orm.Pagination(urlValues)).Select()

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	} else {
		err := database.Handle.Model(&recipeList).Select()

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	return recipeList, nil
}

func (database *DatabaseClient) UpdateRecipe(userRecipe recipes.Recipe) (int, error) {

	var existingRecipe = &recipes.Recipe{UId: userRecipe.UId}

	err := database.Handle.Select(existingRecipe)

	if err == sql.ErrNoRows {
		log.Println(err.Error())
		return 0, err
	}
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	// This user is not the creator of this recipe
	if existingRecipe.UserId != userRecipe.UserId {
		return 0, errors.New("This user is not authorized to update this recipe")
	}
	existingRecipe.Difficulty = userRecipe.Difficulty
	existingRecipe.Name = userRecipe.Name
	existingRecipe.PrepTime = userRecipe.PrepTime
	existingRecipe.Vegetarian = userRecipe.Vegetarian
	existingRecipe.Ratings = userRecipe.Ratings

	err = database.Handle.Update(existingRecipe)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return existingRecipe.UId, nil
}

func (database *DatabaseClient) DeleteRecipe(recipeId int) (int, error) {

	var recipe = &recipes.Recipe{UId: recipeId}
	err := database.Handle.Delete(recipe)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return recipeId, nil
}

// RateRecipe lets users add unlimited rating
func (database *DatabaseClient) RateRecipe(recipeId int, rating int) (int, error) {

	if rating < 1 || rating > 5 {
		return 0, errors.New("Invalid Rating:Must be in the range 1 -5")
	}
	recipe, err := database.GetRecipe(recipeId)
	log.Println(recipe)
	if err == nil && recipe == nil {
		log.Println("Recipe not found")
		return 0, nil
	}

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	recipe.Ratings = append(recipe.Ratings, rating)

	err = database.Handle.Update(recipe)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return rating, nil
}

// SearchRecipes gives paginated search results
func (database *DatabaseClient) SearchRecipes(searchParams map[string]interface{}, urlValues url.Values) ([]recipes.Recipe, error) {
	var searchRecipes = []recipes.Recipe{}
	for key, value := range searchParams {

		if key == "name" {
			var recipesList = []recipes.Recipe{}
			err := database.Handle.Model(&recipesList).Where("name=?", value).Select()
			if err != nil && err != pg.ErrNoRows {
				log.Println(err.Error())
				return nil, err
			}
			searchRecipes = append(searchRecipes, recipesList...)
		}

		if value, ok := value.(bool); ok && key == "vegetarian" {
			var recipesList = []recipes.Recipe{}
			err := database.Handle.Model(&recipesList).Where("vegetarian=?", value).Select()
			if err != nil && err != pg.ErrNoRows {
				log.Println("Error in name")
				return nil, err
			}
			searchRecipes = append(searchRecipes, recipesList...)
		}
		if value, ok := value.(int); ok {
			var recipesList = []recipes.Recipe{}
			if key == "preptime" {
				err := database.Handle.Model(&recipesList).Where("preptime=?", value).Select()
				if err != nil && err != pg.ErrNoRows {
					log.Println(err.Error())
					return nil, err
				}
			} else if strings.EqualFold(key, "difficulty") {
				err := database.Handle.Model(&recipesList).Where("difficulty=?", value).Select()
				if err != nil && err != pg.ErrNoRows {
					log.Println(err.Error())
					return nil, err
				}
			}
			if len(recipesList) > 0 {
				searchRecipes = append(searchRecipes, recipesList...)
			}

		}

	}

	limitValue := urlValues.Get("limit")
	pageValue := urlValues.Get("page")

	var limit int
	var page int
	var err error

	if len(limitValue) > 0 && len(pageValue) > 0 && len(searchRecipes) > 1 {
		limit, err = strconv.Atoi(limitValue)
		log.Println(len(searchRecipes))
		if err != nil {
			log.Println("Failed to parse limit")
			return searchRecipes, nil
		}

		page, err = strconv.Atoi(pageValue)
		if err != nil {
			log.Println("Failed to parse page")
			return searchRecipes, nil
		}

		startIndex := (page - 1) * limit
		if startIndex > -1 && len(searchRecipes) >= startIndex {
			if startIndex+limit <= len(searchRecipes) {
				searchRecipes = searchRecipes[startIndex : startIndex+limit]
				return searchRecipes, nil
			} else {
				searchRecipes = searchRecipes[startIndex:]
				return searchRecipes, nil
			}

		}

	}

	return searchRecipes, nil
}
