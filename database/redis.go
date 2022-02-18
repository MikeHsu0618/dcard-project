package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"sync"
	"time"
)

var ctx = context.Background()
var Redis *redis.Client
var mutex sync.Mutex

var (
	host = os.Getenv("REDIS_HOST")
	port = os.Getenv("REDIS_PORT")
	pws  = os.Getenv("REDIS_PASSWORD")
)

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pws, // no password set
		DB:       0,   // use default DB
	})
}

func Lock(key string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	boolean, err := Redis.SetNX(ctx, key, 1, 3*time.Second).Result()
	if err != nil {
		println(err.Error())
	}
	return boolean
}

func UnLock(key string) int64 {
	nums, err := Redis.Del(ctx, key).Result()
	if err != nil {
		println(err.Error())
		return 0
	}
	return nums
}
