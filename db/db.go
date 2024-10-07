package db

import (
	"fmt"
	"log"
	"strings"

	//_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	//"github.com/jinzhu/gorm"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func ProdEnv() bool {

	environment := os.Getenv("ENVIRONMENT")

	if strings.ToLower(environment) == "development" {
		return false
	}

	return true
}

var db *gorm.DB
var err error

// Init creates a connection to postgress database and
// migrates any new models
func Init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var dbinfo string
	if ProdEnv() {

		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		database := os.Getenv("POSTGRES_DATABASE")

		dbinfo = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=require",
			user,
			password,
			host,
			port,
			database,
		)

	} else {

		user := os.Getenv("TEST_POSTGRES_USER")
		password := os.Getenv("TEST_POSTGRES_PASSWORD")
		host := os.Getenv("TEST_POSTGRES_HOST")
		port := os.Getenv("TEST_POSTGRES_PORT")
		database := os.Getenv("TEST_POSTGRES_DATABASE")

		dbinfo = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
			user,
			password,
			host,
			port,
			database,
		)

	}

	// db, err = gorm.Open("postgres", dbinfo)
	db, err = gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database")
		panic(err)
	}
	log.Println("Database connected")
}

// GetDB ...
func GetDB() *gorm.DB {
	return db
}

// func CloseDB() {
//   db.Close()
// }
