package cache

import (
	"log"
	"time"

	redis "gopkg.in/redis.v4"
)

type CacheClient struct {
	Handle *redis.Client
}

var redisClient = &CacheClient{}

const expirationTime = time.Minute * 30

func GetClient() (*CacheClient, error) {
	err := redisClient.Initialize()

	if err != nil {
		return nil, err
	}

	return redisClient, nil
}

func (redisCache *CacheClient) Initialize() error {

	redisCache.Handle = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	ping, err := redisCache.Handle.Ping().Result()
	if err != nil {
		log.Println(ping, err)
		return err
	}

	return nil
}

func (redisCache *CacheClient) Close() error {
	err := redisCache.Handle.Close()
	if err != nil {
		return err
	}
	return nil
}

func (redisCache *CacheClient) Set(key string, value interface{}) error {

	if value, ok := value.([]byte); ok {
		err := redisCache.Handle.Set(key, value, expirationTime).Err()
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else {
		err := redisCache.Handle.Set(key, value, expirationTime).Err()
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	return nil
}

func (redisCache *CacheClient) Get(key string) (string, error) {
	result := redisCache.Handle.Get(key)
	if result.Err() != nil {
		return "", result.Err()
	}

	value := result.Val()
	return value, nil
}

func (redisCache *CacheClient) Delete(key string) (bool, error) {
	result := redisCache.Handle.Del(key)
	if result.Err() != nil {
		return false, result.Err()
	}

	return true, nil
}
