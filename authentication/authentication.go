package authentication

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	cache "recipes/cache"
	auth "recipes/models/authentication"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
)

func CreateToken(username string, userId int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  username,
		"userid":    userId,
		"timestamp": time.Now(),
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		log.Println("Error signing token")
		panic(err.Error)
	}

	cacheClient, err := cache.GetCacheClient(cache.REDIS)
	if err != nil {
		log.Println("Failed to get cache client")
	}

	defer cacheClient.Close()

	err = cacheClient.Set(username, tokenString)
	if err != nil {
		log.Println("Failed to set token in the cache")
	}

	return tokenString
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(map[string]string{"error": error.Error()})
					return
				}
				if token.Valid {

					cacheClient, err := cache.GetCacheClient(cache.REDIS)
					if err != nil {
						log.Println("Failed to get cache client")
						json.NewEncoder(w).Encode(map[string]string{"error": error.Error()})
						return
					}
					defer cacheClient.Close()
					claims := token.Claims
					var tokenClaim auth.TokenClaims
					err = mapstructure.Decode(claims.(jwt.MapClaims), &tokenClaim)
					if err != nil {
						log.Println("Decode error")
						log.Print(err.Error())
						json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
						return
					}
					_, err = cacheClient.Get(tokenClaim.Username)
					if err != nil {
						log.Println("Failed to set token in the cache")
						json.NewEncoder(w).Encode(map[string]string{"error": "Token has expired. Please authorize again."})
						return
					}
					context.Set(req, "decoded", tokenClaim)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(map[string]string{"error": "AUthorization token is required"})
		}
	})
}
