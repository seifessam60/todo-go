package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB



func ConnectDB(){
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_PORT"),
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_NAME"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to Connect to database: ", err)
	}
	log.Println("Database Connected Successfully!")
}

func Migrate(models ...interface{}) {
	if err := DB.AutoMigrate(models...); err != nil{
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Database migrated successfully!")

}