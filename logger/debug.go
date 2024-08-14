package logger

import (
    "log"
	"strings"
	"os"
)


func LogTextSpace(text string){
	environment := os.Getenv("ENVIRONMENT")
	
	if strings.ToLower(environment) == "development" {
		log.Println("")	
		log.Println(text)	
		log.Println("")		
	}
	
}

func LogHeaderSpace(text string){
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

func LogTextValueSpace(text string, value interface{}){

	environment := os.Getenv("ENVIRONMENT")

	if strings.ToLower(environment) == "development" {
		log.Println("")	
		log.Println(text, value)	
		log.Println("")		
	}

}


func LogTextValue(text string, value interface{}){

	environment := os.Getenv("ENVIRONMENT")

	if strings.ToLower(environment) == "development" {
		log.Println(text, value)		
	}

}