package repository

import (
	"sync"

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

func NewUrlRepo(db *gorm.DB, client *redis.Client, logger *logger.Logger) UrlRepository {
	return &UrlRepo{
		db:     db,
		redis:  client,
		logger: logger,
	}
}

func (r *UrlRepo) Create(url *Url) (err error) {
	if err := r.db.Create(url).Error; err != nil {
		return err
	}
	return nil
}

func (r *UrlRepo) GetById(urlId int64, url *Url) (err error) {
	if err := r.db.Where("id", urlId).First(&url).Error; err != nil {
		return err
	}
	return nil
}

func (r *UrlRepo) GetByOrgUrl(orgUrl string) (url *Url, err error) {
	if err := r.db.Where("org_url", orgUrl).First(&url).Error; err != nil {
		return url, err
	}
	return url, nil
}
