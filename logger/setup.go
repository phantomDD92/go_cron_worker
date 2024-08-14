package logger

import (
	"encoding/json"
	"fmt"
	//"io"
	"os"
	"strings"
	"go_proxy_worker/utils"
	log "github.com/sirupsen/logrus"
)

// JSONFormatter is a logger for use with Logrus
type JSONFormatter struct {
	Program string
	Env     string
}

type LogData struct {
	File  string
	Error error
}

func UseJSONLogFormat() {
	/*
		ENABLE LOGGER
	*/
	
	env := utils.GetEnv("GIN_ENV", "production")
	program := utils.GetEnv("SERVICE_NAME", "go_proxy_worker")

	// var f *os.File

	// filename := "log_file.log"
	// if _, err := os.Stat(filename); os.IsNotExist(err) {
	// 	var fileError error
	// 	f, fileError = os.Create(filename)
	// 	if fileError != nil {
	// 		log.WithFields(log.Fields{
	// 			"file": "utils.log",
	// 		}).Error("issue creating log file")
	// 	}
	// } else {
	// 	var fileError error
	// 	f, fileError = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0660)
	// 	if fileError != nil {
	// 		log.WithFields(log.Fields{
	// 			"file": "utils.log",
	// 		}).Error("issue creating log file")
	// 	}
	// }

	// mw := io.MultiWriter(os.Stdout, f)
	// log.SetOutput(mw)

	log.SetOutput(os.Stdout)

	log.SetFormatter(&JSONFormatter{
		Program: program,
		Env:     env,
	})

  }


  // Timestamps in microsecond resolution (like time.RFC3339Nano but microseconds)
var timeStampFormat = "2006-01-02T15:04:05.000000Z07:00"

// Format includes the program, environment, and a custom time format: microsecond resolution
func (f *JSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	data := make(log.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		data[k] = v
	}
	data["time"] = entry.Time.UTC().Format(timeStampFormat)
	data["msg"] = entry.Message
	data["level"] = strings.ToUpper(entry.Level.String())
	data["program"] = f.Program
	data["env"] = f.Env

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}