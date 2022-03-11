package database

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"os"
)

var Db *gorm.DB

func init() {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	db := os.Getenv("POSTGRES_DB")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v "+
		"sslmode=disable TimeZone=Asia/Shanghai", host, user, pwd, db, port)

	hostSlave := os.Getenv("POSTGRES_SLAVE_HOST")
	userSlave := os.Getenv("POSTGRES_SLAVE_USER")
	dbSlave := os.Getenv("POSTGRES_SLAVE_DB")
	pwdSlave := os.Getenv("POSTGRES_SLAVE_PASSWORD")
	portSlave := os.Getenv("POSTGRES_SLAVE_PORT")
	dsnSlave := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v "+
		"sslmode=disable TimeZone=Asia/Shanghai", hostSlave, userSlave, pwdSlave, dbSlave, portSlave)

	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("postgres connect error %v", err)
	}

	if Db.Error != nil {
		fmt.Printf("database error %v", Db.Error)
	}

	err = Db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{postgres.Open(dsnSlave)},
	}))
	if err != nil {
		fmt.Printf("dbresolver connect error %v", err)
	}
	if Db.Error != nil {
		fmt.Printf("database error %v", Db.Error)
	}
}
