package cache

import (
	redis "recipes/cache/redis"
)

func GetCacheClient(clientName string) (CacheClient, error) {
	var client CacheClient
	var err error
	switch clientName {
	case REDIS:
		client = &redis.CacheClient{}
		err := client.Initialize()

		if err != nil {
			return nil, err
		}
	}

	return client, err
}
