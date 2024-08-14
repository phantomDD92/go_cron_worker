package workers

import (
	"github.com/go-co-op/gocron"
	"go_proxy_worker/logger"
	"go_proxy_worker/utils"
	// "go_proxy_worker/dbRedisQueries"
	// "github.com/fatih/structs"
	"go_proxy_worker/db"
	// "github.com/jinzhu/gorm"
	// "strconv"
	// "strings"
	"time"
	"log"
	// "fmt"
	"encoding/json"
)


type FailedResponse struct {
    TimeWindow 				string `json:"time_window"`
	Method					string `json:"method"`

    Proxy				    string `json:"proxy"`
    ProxyUrl                string `json:"proxy_url"`
    ProxyStatusCode			int `json:"proxy_status_code"`
    ProxyApiCredits			int64 `json:"proxy_api_credits"`
	ProxyUniqueId           string `json:"proxy_unique_id"`

    BodyString				string `json:"body_string"`
    ContentType             string `json:"content_type"`
	
	Url   					string `json:"url"`
	Domain					string `json:"domain"`

    FailedReason			string `json:"failed_reason"`
	BlockType				string `json:"block_type"`
	Block					string `json:"block"`

}



func CronLogFailedValidationResponses() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(15).Minutes().Do(LogFailedValidationResponses)
	s.StartBlocking()
}

func LogFailedValidationResponses(){


	/*
		Pulls responses that returned 200 but failed validation checks from
		queue, and logs meta data in DB and response in bucket.
	*/

	


	log.Println("Run LogFailedValidationResponses")

	fileName := "cron_log_failed_validation_responses.go"

	// Redis Details
	var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	// var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()


	queueName := utils.RedisEnvironmentVersion("failedResponsesQueue")

	failedResponseKeys, err := coreProxyRedisClient.SMembers(redisContext, queueName).Result()
	if err != nil {
		errData := map[string]interface{}{
			"queueName": queueName,
		}
		logger.LogError("INFO", fileName, err, "failedResponsesQueue not in Redis", errData)
	}


	log.Println("failedResponseKeys", failedResponseKeys)

	for _, failedResponseKey := range failedResponseKeys {

		// Delete From Queue
		err = coreProxyRedisClient.SRem(redisContext, queueName, failedResponseKey).Err()
		if err != nil {
			errData := map[string]interface{}{
				"failedResponseKey": failedResponseKey,
				"queueName": queueName,
			}
			logger.LogError("ERROR", fileName, err, "failed to delete failedResponseKey from failedResponsesQueue in Redis", errData)
		}

		// Get Failed Response
		failedResponseString, err := coreProxyRedisClient.Get(redisContext, failedResponseKey).Result()
		if err != nil {
			errData := map[string]interface{}{
				"failedResponseKey": failedResponseKey,
			}
			logger.LogError("WARN", fileName, err, "failedResponseKey not in Redis", errData)
		}


		/*
			If Failed Response Is In Redis

			- Store the meta data in db
			- Store the HTML response in bucket
		*/

		if failedResponseString != "" {

			var failedResponse FailedResponse
			json.Unmarshal([]byte(failedResponseString), &failedResponse)


			/*
				Store HTML Response In Bucket
			*/

			// Convert HTML String to HTML File


			// Create File Path
			// filePath := 

			// Store In File in DO Bucket
			// UploadSpacesBucket(spaces string, filePath string, file multipart.File)

			// Insert Meta Data Into DB


		}



		// Delete Failed Response From Redis
		_, err = coreProxyRedisClient.Del(redisContext, failedResponseKey).Result()
		if err != nil {
			errData := map[string]interface{}{
				"failedResponseKey": failedResponseKey,
			}
			logger.LogError("WARN", fileName, err, "error deleting failedResponseKey in Redis", errData)
		}



	}


}