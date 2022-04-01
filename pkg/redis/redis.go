package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
)

var client *redis.Client

type Config struct {
	Address string
	Pws     string
}

func NewRedisClient(config Config) *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Pws, // no password set
		DB:       0,          // use default DB
	})

	return client
}

func NewTestRedisClient() *redis.Client {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client = redis.NewClient(&redis.Options{
		Addr: s.Addr(), // mock redis server的地址
	})
	return client
}

func GetRedisClient() *redis.Client {
	return client
}
