package cache

const (
	REDIS string = "redis"
)

type CacheClient interface {
	Initialize() error
	Close() error
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Delete(key string) (bool, error)
}
