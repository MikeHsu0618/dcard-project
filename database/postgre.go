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
	user := os.Getenv("POSTGRES_USER")
	db := os.Getenv("POSTGRES_DB")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("POSTGRES_PORT")

	var err error
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%v "+
		"sslmode=disable TimeZone=Asia/Shanghai", user, db, pwd, port)
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}

	if Db.Error != nil {
		fmt.Printf("database error %v", Db.Error)
	}
}
