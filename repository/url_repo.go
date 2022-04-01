package repository

import (
	"context"
	"strconv"
	"sync"
	"time"

	"dcard-project/model"
	"dcard-project/pkg/logger"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UrlRepo struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *logger.Logger
}

var mutex sync.Mutex

func NewUrlRepo(db *gorm.DB, client *redis.Client, logger *logger.Logger) model.UrlRepository {
	return &UrlRepo{
		db:     db,
		redis:  client,
		logger: logger,
	}
}

func (r *UrlRepo) Create(url *model.Url) (err error) {
	if err := r.db.Create(url).Error; err != nil {
		return err
	}
	return nil
}

func (r *UrlRepo) GetById(urlId int64, url *model.Url) (err error) {
	if err := r.db.Where("id", urlId).First(&url).Error; err != nil {
		return err
	}
	return nil
}

func (r *UrlRepo) GetByOrgUrl(orgUrl string) (url *model.Url, err error) {
	if err := r.db.Where("org_url", orgUrl).First(&url).Error; err != nil {
		return url, err
	}
	return url, nil
}

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
