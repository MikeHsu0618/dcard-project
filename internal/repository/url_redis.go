package repository

import (
	"context"
	"strconv"
	"time"
)

func (r *UrlRepo) Lock(key string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	boolean, err := r.redis.SetNX(context.Background(), key, 1, 5*time.Second).Result()
	if err != nil {
		r.logger.Error("redis:" + err.Error())
		return true
	}
	return boolean
}

func (r *UrlRepo) UnLock(key string) int64 {
	nums, err := r.redis.Del(context.Background(), key).Result()
	if err != nil {
		r.logger.Error("redis:" + err.Error())
		return 0
	}
	return nums
}

func (r *UrlRepo) GetCache(key int64) (result string, err error) {
	return r.redis.Get(context.Background(), strconv.FormatInt(key, 10)).Result()
}

func (r *UrlRepo) SetCache(key int64, orgUrl string) {
	_, err := r.redis.Set(
		context.Background(),
		strconv.FormatInt(key, 10),
		orgUrl,
		30*24*time.Hour).Result()
	if err != nil {
		println(err.Error())
	}
}
