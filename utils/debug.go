package utils

import (
	"log"
	"os"
	"strings"
)

func OnlyRunTestAccounts() (bool, []string) {
	args := os.Args[1:]
	if len(args) == 2 {
		if args[1] == "test_accounts" {
			testAccounts := []string{"1103", "3"}
			return true, testAccounts
		}
	}

	emptyList := []string{}
	return false, emptyList
}

func LogTextSpace(text string) {
	environment := os.Getenv("ENVIRONMENT")

	if strings.ToLower(environment) == "development" {
		log.Println("")
		log.Println(text)
		log.Println("")
	}

}

func LogHeaderSpace(text string) {
	environment := os.Getenv("ENVIRONMENT")

	if strings.ToLower(environment) == "development" {
		log.Println("")
		log.Println("#################")
		log.Println("")
		log.Println(strings.ToUpper(text))
		log.Println("")
		log.Println("#################")
		log.Println("")
	}

}

func LogTextValue(text string, value interface{}) {

	environment := os.Getenv("ENVIRONMENT")

	if strings.ToLower(environment) == "development" {
		log.Println("")
		log.Println(text, value)
		log.Println("")
	}

}
