package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

var client *redis.Client

type Config struct {
	address string
	pws     string
}

func NewRedisClient() *redis.Client {
	config := Config{
		address: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		pws:     os.Getenv("REDIS_PASSWORD"),
	}

	client = redis.NewClient(&redis.Options{
		Addr:     config.address,
		Password: config.pws, // no password set
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
