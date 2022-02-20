package database

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}

	if Db.Error != nil {
		fmt.Printf("database error %v", Db.Error)
	}
}
