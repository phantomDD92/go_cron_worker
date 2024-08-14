package logger

import (
	"os"
	"strings"
	log "github.com/sirupsen/logrus"
)

func LogError(errorType string, fileName string, err error, msg string, errData map[string]interface{}) {

	
	
	var logFields log.Fields
	logFields = log.Fields{
		"file":               fileName,	
		"error":              err,
		"error_data": 		  errData,
	}
	

	if errorType == "DEBUG" {
		
		environment := os.Getenv("ENVIRONMENT")
		if strings.ToLower(environment) != "development" {
			log.WithFields(logFields).Debug(msg)
		}

	} else if errorType == "INFO" {

		environment := os.Getenv("ENVIRONMENT")
		if strings.ToLower(environment) != "development" {
			log.WithFields(logFields).Info(msg)
		}

	} else if errorType == "WARN" {

		environment := os.Getenv("ENVIRONMENT")
		if strings.ToLower(environment) != "development" {
			log.WithFields(logFields).Info(msg)
		}

	} else if errorType == "ERROR" {
		
		log.WithFields(logFields).Error(msg)

	} else if errorType == "FATAL" {

		log.WithFields(logFields).Fatal(msg)
	}

	

	

}
