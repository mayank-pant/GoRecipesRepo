package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	authentication "recipes/authentication"
	recipehandler "recipes/handlers/recipe"
	userhandler "recipes/handlers/user"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	redis "gopkg.in/redis.v4"
)

func initializeDatabases() {

	DB_HOST := os.Getenv("DB_HOST")
	fmt.Println("HOST")
	fmt.Println(DB_HOST)
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, "kapeiq", "kapeiq", "kapeiq")

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("PING ERR")
		return
	}

	db.Query("DROP TABLE users,passwords,recipes;")
	db.Query("CREATE TABLE users" +
		"(" +
		"uid serial NOT NULL," +
		"username character varying(200) NOT NULL," +
		"isactive boolean NOT NULL," +
		"created date default now()," +
		"CONSTRAINT users_pkey PRIMARY KEY (uid)" +
		")WITH (OIDS=FALSE);")

	db.Query("CREATE TABLE passwords (" +
		"user_id serial NOT NULL," +
		"Hash character varying(200) NOT NULL," +
		"created TIMESTAMP NOT NULL DEFAULT NOW()," +
		"PRIMARY KEY(user_id)," +
		"CONSTRAINT userid_fkey FOREIGN KEY (user_id) REFERENCES USERS (uid)" +
		")WITH(OIDS=FALSE);")

	db.Query("CREATE TABLE recipes" +
		"(" +
		"uid serial NOT NULL," +
		"name character varying(200) NOT NULL," +
		"vegetarian boolean NOT NULL," +
		"prep_time integer NOT NULL," +
		"difficulty integer NOT NULL CHECK(difficulty BETWEEN 1 AND 3)," +
		"user_id integer NOT NULL," +
		"created date default now()," +
		"ratings integer[]," +
		"CONSTRAINT recipes_pkey PRIMARY KEY (uid)," +
		"CONSTRAINT userid_fkey FOREIGN KEY (user_id) REFERENCES USERS(uid)" +
		")WITH (OIDS=FALSE);")
}

func pingCheck() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	ping, err := client.Ping().Result()
	if err != nil {
		log.Println("Redis failed to start")
		panic(err)
	}
	log.Println(ping, err)
}

var router = mux.NewRouter().StrictSlash(true)

func main() {
	pingCheck()
	log.Println("Starting the recipe app ...")
	router.HandleFunc("/signup", userhandler.SignUp).Methods("POST", "OPTIONS")
	router.HandleFunc("/authenticate", userhandler.ObtainToken).Methods("POST", "OPTIONS")
	router.HandleFunc("/recipes", recipehandler.ListRecipes).Methods("GET", "OPTIONS")
	router.HandleFunc("/recipes", authentication.ValidateMiddleware(recipehandler.CreateRecipe)).Methods("POST")
	router.HandleFunc("/recipes/{id}", recipehandler.GetRecipe).Methods("GET")
	router.HandleFunc("/recipes/{id}", authentication.ValidateMiddleware(recipehandler.UpdateRecipe)).Methods("PUT")
	router.HandleFunc("/recipes/{id}", authentication.ValidateMiddleware(recipehandler.DeleteRecipe)).Methods("DELETE")
	router.HandleFunc("/recipes/{id}/rating", recipehandler.RateRecipe).Methods("POST")
	router.HandleFunc("/search/recipes", recipehandler.SearchRecipes).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
}

func init() {
	log.Println("Initializing Databases")
	initializeDatabases()
}
