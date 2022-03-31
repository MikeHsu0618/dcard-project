package postgres

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var db *gorm.DB
var err error

type DbConf struct {
	host string
	user string
	db   string
	pwd  string
	port string
}

func NewPgClient() *gorm.DB {
	master := DbConf{
		host: os.Getenv("POSTGRES_HOST"),
		user: os.Getenv("POSTGRES_USER"),
		db:   os.Getenv("POSTGRES_DB"),
		pwd:  os.Getenv("POSTGRES_PASSWORD"),
		port: os.Getenv("POSTGRES_PORT"),
	}

	slave := DbConf{
		host: os.Getenv("POSTGRES_SLAVE_HOST"),
		user: os.Getenv("POSTGRES_SLAVE_USER"),
		db:   os.Getenv("POSTGRES_SLAVE_DB"),
		pwd:  os.Getenv("POSTGRES_SLAVE_PASSWORD"),
		port: os.Getenv("POSTGRES_SLAVE_PORT"),
	}

	db, err = gorm.Open(postgres.Open(GetPgDns(master)), &gorm.Config{})
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{postgres.Open(GetPgDns(slave))},
	}))

	if err != nil || db.Error != nil {
		panic("database error")
	}

	return db
}

func GetPgDns(conf DbConf) (dsn string) {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Taipei",
		conf.host, conf.user, conf.pwd, conf.db, conf.port)
}
