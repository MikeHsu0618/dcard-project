package main

import (
	"os"
	"testing"

	"dcard-project/internal/repository"
	"dcard-project/internal/service"
	"dcard-project/pkg/logger"
	rds "dcard-project/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	rdsClient *redis.Client
	lg        *logger.Logger
	svc       repository.UrlService
	r         = gin.Default()
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardrop()
	os.Exit(code)
}

func teardrop() {
	initDBTable()
	println("tearDrop")
}

func setup() {
	println("setup")
	initService()
}

func initService() {
	db, _ = gorm.Open(sqlite.Open("./url_test.db"), &gorm.Config{})
	rdsClient = rds.NewTestRedisClient()
	lg = logger.NewLogger()

	repo := repository.NewUrlRepo(db, rdsClient, lg)
	svc = service.NewUrlService(repo)
}

func initDBTable() {
	var urls []repository.Url
	db.Find(&urls)
	db.Delete(&urls)
}
