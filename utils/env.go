package utils

import (
	"os"
	"fmt"
	"strings"
	log "github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
)

func LoadEnv(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetEnv returns an environment variable or a default value if not present
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}


// LoadEnvVars will load a ".env[.development|.test]" file if it exists and set ENV vars.
// Useful in development and test modes. Not used in production.
func LoadEnvVars() {
	env := GetEnv("GIN_ENV", "development")

	if env == "production" || env == "staging" {
		log.Println("Not using .env file in production or staging.")
		return
	}

	filename := ".env." + env

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = ".env"
	}

	err := godotenv.Load(filename)
	if err != nil {
		log.Warn(".env file not loaded")
	}
}


func RedisEnvironmentVersion(redisString string) string {

	environment := os.Getenv("ENVIRONMENT")
	
	if strings.ToLower(environment) == "development" {
		redisString = "DEV_" + fmt.Sprintf("%v", redisString) 		
	}

	return redisString
}


func ProdEnv() bool {

	environment := os.Getenv("ENVIRONMENT")
	
	if strings.ToLower(environment) == "development" {
		return false		
	}

	return true
}