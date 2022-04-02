package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var db *gorm.DB
var err error

type Config struct {
	Host string
	User string
	Db   string
	Pwd  string
	Port string
}

func NewPgClient(master Config, slave Config) *gorm.DB {
	db, err = gorm.Open(postgres.Open(GetPgDns(master)), &gorm.Config{})
	if err != nil || db.Error != nil {
		panic("master database error")
	}

	if len(slave.Host) > 0 {
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{postgres.Open(GetPgDns(slave))},
		}))

		if err != nil || db.Error != nil {
			panic("slave database error")
		}
	}

	return db
}

func GetPgDns(conf Config) (dsn string) {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Taipei",
		conf.Host, conf.User, conf.Pwd, conf.Db, conf.Port)
}
