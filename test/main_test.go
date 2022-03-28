package main

import (
    "dcard-project/model"
    "dcard-project/pkg/migration"
    "dcard-project/pkg/postgres"
    rds "dcard-project/pkg/redis"
    "dcard-project/pkg/testutil"
    "dcard-project/repository"
    "dcard-project/service"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "os"
    "path"
    "runtime"
    "strings"
    "testing"
)

const (
    file = ".env.testing"
)

var (
    db     *gorm.DB
    client *redis.Client
    svc    model.UrlService
    r      = gin.Default()
)

func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    teardrop()
    os.Exit(code)
}

func teardrop() {
    println("tearDrop")
}

func setup() {
    println("setup")
    testutil.LoadEnv(file)
    initService()
}

func initService() {
    db = postgres.NewPgClient()
    client = rds.NewRedisClient()
    repo := repository.NewUrlRepo(db, client)
    svc = service.NewUrlService(repo)
}

func initDBTable() {
    m := migration.NewMigrate(migration.Config{
        DatabaseDriver: migration.PostgresDriver,
        DatabaseURL:    `postgres://postgres:postgres@localhost:5432/` + os.Getenv("POSTGRES_DB") + `?sslmode=disable`,
        SourceDriver:   migration.FileDriver,
        SourceURL:      getMigrationPathByCaller(),
        SchemaName:     "schema_migrations",
    })
    err := m.Refresh()
    if err != nil {
        println("reset database Error")
    }
}

func getMigrationPathByCaller() string {
    var abPath string
    _, filename, _, ok := runtime.Caller(0)
    if ok {
        abPath = path.Dir(filename)
    }
    split := strings.Split(abPath, "/")
    split[len(split)-1] = "migrations"
    mPath := strings.Join(split, "/")
    return mPath
}
