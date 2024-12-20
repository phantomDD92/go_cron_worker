package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"

	"github.com/fatih/structs"
	"github.com/go-co-op/gocron"

	// "go_proxy_worker/models"
	"go_proxy_worker/logger"
	"go_proxy_worker/models"
	"go_proxy_worker/slack"
	"go_proxy_worker/utils"

	"github.com/go-redis/redis/v8"

	// "go_proxy_worker/dbRedisQueries"
	// "github.com/fatih/structs"
	"go_proxy_worker/db"
	// // "gorm.io/gorm"
	// "strconv"
	// "strings"

	"time"
	// "fmt"
	// "os"
)

type SopsProxyTestResponse struct {
	UserProxyTestResults []map[string]interface{} `json:"user_proxy_test_results"`
	TestResults          []map[string]interface{} `json:"test_results"`
	WorkingProxies       []WorkingProxy           `json:"working_proxies"`
	CategoryProxies      []map[string]interface{} `json:"category_proxies"`
}

type WorkingProxy struct {
	Proxy                 string            `json:"proxy"`
	ConcurrencyLimit      int               `json:"concurrency_limit"`
	Type                  string            `json:"type"`
	Features              map[string]string `json:"features"`
	TotalRequests         int               `json:"total_requests"`
	Successful            int               `json:"successful"`
	Test1Passed           int               `json:"test_1_passed"`
	TotalLatency          int               `json:"total_latency"`
	SuccessRate           float64           `json:"success_rate"`
	ValidRate             float64           `json:"valid_rate"`
	AvgSuccessLatency     float64           `json:"avg_success_latency"`
	CPM                   int               `json:"CPM"`
	PremiumDomain         bool              `json:"premium_domain"`
	PremiumDomainMultiple int               `json:"premium_domain_multiple"`
}

func CronProxyTesterQueue() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Seconds().Do(RunProxyTesterQueue)
	s.StartBlocking()
}

func RunProxyTesterQueue() {

	// if utils.OnlyRunTestAccounts() {
	// 	log.Println("helloe")
	// }

	fileName := "cron_run_proxy_tester_queue.go"
	emptyErrMap := make(map[string]interface{})

	// Redis Details
	var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	// // Getting the length of the list
	// listLength, err := coreProxyRedisClient.LLen(redisContext, "proxyTestsQueue").Result()
	// if err != nil {
	// 	logger.LogError("ERROR", fileName, err, "error getting length of proxyTestsQueue Redis", emptyErrMap)
	// }

	var wg sync.WaitGroup

	// Loop through and pop each element until the list is empty
	for {
		value, err := coreProxyRedisClient.LPop(redisContext, "proxyTestsQueue").Result()
		if err == redis.Nil {
			// The list is empty, break the loop
			break
		} else if err != nil {
			logger.LogError("ERROR", fileName, err, "error getting element from proxyTestsQueue Redis", emptyErrMap)
			continue
		}

		// Convert the popped value (JSON string) into a map
		var proxyTestSetup map[string]interface{}
		if err := json.Unmarshal([]byte(value), &proxyTestSetup); err != nil {
			logger.LogError("ERROR", fileName, err, "Error unmarshalling JSON:", emptyErrMap)
			continue
		}

		wg.Add(1)
		go MakeProxyTesterRequest(&wg, proxyTestSetup, fileName)

		// if proxyTestSetup["test_type"] == nil || proxyTestSetup["test_type"] == "test_internal_proxy_pools" || proxyTestSetup["test_type"] == "" {
		// 	wg.Add(1)
		// 	go MakeInternalProxyTesterRequest(&wg, proxyTestSetup, fileName)
		// } else {
		// 	wg.Add(1)
		// 	go MakeInternalProxyTesterRequest(&wg, proxyTestSetup, fileName)
		// 	wg.Add(1)
		// 	go MakeUserProxyTesterRequest(&wg, proxyTestSetup, fileName)
		// }

	}

	wg.Wait()

}

func MakeProxyTesterRequest(wg *sync.WaitGroup, proxyTestSetup map[string]interface{}, fileName string) {

	defer wg.Done()

	emptyErrMap := make(map[string]interface{})

	// load DB
	var db = db.GetDB()

	logger.LogTextValue("proxyTestSetup", proxyTestSetup)

	client := &http.Client{
		Timeout: 130 * time.Second,
	}

	// Create Request
	var req *http.Request
	var proxyTestEndpoint string

	// proxyTestSetup["api_key"] = "85bb39cd-c5a6-44cc-9221-401e860e52b1"
	postBody, _ := json.Marshal(proxyTestSetup)
	postBodyBytes := bytes.NewBuffer(postBody)

	if proxyTestSetup["test_type"] == nil || proxyTestSetup["test_type"] == "test_internal_proxy_pools" || proxyTestSetup["test_type"] == "" {
		// Create Proxy Provider Test Request
		proxyTestEndpoint = "https://backend.scrapeops.io/test-proxies-providers/v2/"
		req, _ = http.NewRequest("POST", proxyTestEndpoint, postBodyBytes)
		req.Header.Set("Content-Type", "application/json")
	} else {
		// test_user_proxy_settings
		// Create User Proxy Settings Test Request
		proxyTestEndpoint = "https://backend.scrapeops.io/test-user-proxy-settings/v1/"
		req, _ = http.NewRequest("POST", proxyTestEndpoint, postBodyBytes)
		req.Header.Set("Content-Type", "application/json")
	}

	logger.LogTextValue("proxyTestEndpoint", proxyTestEndpoint)
	logger.LogTextValue("Sent Request", proxyTestSetup["test_id"])

	// Update DB With Test Status
	sopsProxyTestResultMap := map[string]interface{}{
		"test_status": "processing",
	}

	var sopsProxyTestResult models.SopsProxyTestResult
	result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
	if result.Error != nil || result.RowsAffected == 0 {
		errData := structs.Map(sopsProxyTestResult)
		logger.LogError("ERROR", fileName, result.Error, "Failed to update sopsProxyTestResultMap in DB", errData)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.LogTextValue("error", err)
		logger.LogTextValue("proxyTestSetup", proxyTestSetup)

		// Update DB for Failed Test
		sopsProxyTestResultMap["test_status"] = "failed"
		result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
		if result.Error != nil || result.RowsAffected == 0 {
			errData := structs.Map(sopsProxyTestResult)
			logger.LogError("ERROR", fileName, result.Error, "Failed to update or create sopsProxyTestResultMap in DB", errData)
		}

	} else {
		logger.LogTextValue("Test Response", "")
		logger.LogTextValueSpace("proxy response code", resp.StatusCode)

		defer resp.Body.Close()

		// Convert Body To JSON
		var sopsProxyTestResponse SopsProxyTestResponse
		json.NewDecoder(resp.Body).Decode(&sopsProxyTestResponse)

		// Convert the struct to JSON
		jsonSopsProxyTestResponse, err := json.Marshal(sopsProxyTestResponse)
		if err != nil {
			logger.LogError("ERROR", fileName, result.Error, "Error converting sopsProxyTestResponse to JSON", emptyErrMap)
		}

		sopsProxyTestResultMap["test_status"] = "completed"
		sopsProxyTestResultMap["test_results"] = jsonSopsProxyTestResponse

		logger.LogTextValue("", "")
		logger.LogTextValue("sopsProxyTestResponse", sopsProxyTestResponse)
		// logger.LogTextValue("sopsProxyTestResultMap", sopsProxyTestResultMap)

		// Update DB for Failed Test
		result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
		if result.Error != nil || result.RowsAffected == 0 {
			errData := structs.Map(sopsProxyTestResult)
			logger.LogError("ERROR", fileName, result.Error, "Failed to update or create accountProxyStats in DB", errData)
		}

		/*

			SEND SLACK MESSAGE

		*/
		if proxyTestSetup["test_type"] == nil || proxyTestSetup["test_type"] == "test_internal_proxy_pools" || proxyTestSetup["test_type"] == "" {

			// Create Stats Blob
			statsString := "*URL:* " + fmt.Sprintf("%v", proxyTestSetup["test_url"]) + "\n"
			statsString = statsString + "*Test Mode:* " + fmt.Sprintf("%v", proxyTestSetup["test_mode"]) + "\n"
			statsString = statsString + "*Num Requests:* " + fmt.Sprintf("%v", proxyTestSetup["num_test_requests"]) + "\n\n"
			statsString = statsString + "```Proxy | Success | Valid | Latency |\n"
			for _, workingProxy := range sopsProxyTestResponse.WorkingProxies {
				statsString = statsString + workingProxy.Proxy + " | "
				statsString = statsString + fmt.Sprintf("%v", int(workingProxy.SuccessRate)) + "% | "
				statsString = statsString + fmt.Sprintf("%v", int(workingProxy.ValidRate)) + "% | "
				statsString = statsString + fmt.Sprintf("%v", (math.Round(workingProxy.AvgSuccessLatency*100)/100)) + "s |\n"

			}
			statsString = statsString + "```"

			// Send Slack Message
			headline := "Proxy Test: " + fmt.Sprintf("%v", proxyTestSetup["test_domain"])
			slack.SlackStatsAlert("#proxy-optimizer-tests", headline, statsString)

		}

	}

	return
}

func MakeInternalProxyTesterRequest(wg *sync.WaitGroup, proxyTestSetup map[string]interface{}, fileName string) {

	defer wg.Done()

	emptyErrMap := make(map[string]interface{})

	// load DB
	var db = db.GetDB()

	logger.LogTextValue("proxyTestSetup", proxyTestSetup)

	client := &http.Client{
		Timeout: 130 * time.Second,
	}

	// Create Request
	var req *http.Request
	var proxyTestEndpoint string

	// proxyTestSetup["api_key"] = "85bb39cd-c5a6-44cc-9221-401e860e52b1"
	postBody, _ := json.Marshal(proxyTestSetup)
	postBodyBytes := bytes.NewBuffer(postBody)

	// Create Proxy Provider Test Request
	proxyTestEndpoint = "https://backend.scrapeops.io/test-proxies-providers/v2/"
	req, _ = http.NewRequest("POST", proxyTestEndpoint, postBodyBytes)
	req.Header.Set("Content-Type", "application/json")

	logger.LogTextValue("Internal Test proxyTestEndpoint", proxyTestEndpoint)
	logger.LogTextValue("Internal Test Sent Request", proxyTestSetup["test_id"])

	var sopsProxyTestResult models.SopsProxyTestResult

	// Update DB With Test Status
	sopsProxyTestResultMap := map[string]interface{}{}

	if proxyTestSetup["test_type"] == "test_internal_proxy_pools" {
		sopsProxyTestResultMap["test_status"] = "processing"
		result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
		if result.Error != nil || result.RowsAffected == 0 {
			errData := structs.Map(sopsProxyTestResult)
			logger.LogError("ERROR", fileName, result.Error, "Failed to update sopsProxyTestResultMap in DB", errData)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.LogTextValue("Internal Test error", err)
		logger.LogTextValue("Internal Test proxyTestSetup", proxyTestSetup)

		// Update DB for Failed Test
		if proxyTestSetup["test_type"] == "test_internal_proxy_pools" {
			sopsProxyTestResultMap["test_status"] = "failed"
			result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
			if result.Error != nil || result.RowsAffected == 0 {
				errData := structs.Map(sopsProxyTestResult)
				logger.LogError("ERROR", fileName, result.Error, "Failed to update or create sopsProxyTestResultMap in DB", errData)
			}
		}

	} else {
		logger.LogTextValue("Internal Test Response", "")
		logger.LogTextValueSpace("Internal Test response code", resp.StatusCode)

		defer resp.Body.Close()

		// Convert Body To JSON
		var sopsProxyTestResponse SopsProxyTestResponse
		json.NewDecoder(resp.Body).Decode(&sopsProxyTestResponse)

		// logger.LogTextValue("", "")
		logger.LogTextValue("Internal Test sopsProxyTestResponse", sopsProxyTestResponse)
		// logger.LogTextValue("sopsProxyTestResultMap", sopsProxyTestResultMap)

		if proxyTestSetup["test_type"] == "test_internal_proxy_pools" {
			sopsProxyTestResultMap["test_status"] = "completed"

			// Convert the struct to JSON
			jsonSopsProxyTestResponse, err := json.Marshal(sopsProxyTestResponse)
			if err != nil {
				logger.LogError("ERROR", fileName, err, "Error converting sopsProxyTestResponse to JSON", emptyErrMap)
			}
			sopsProxyTestResultMap["test_results"] = jsonSopsProxyTestResponse

			// Update DB for Successful Test
			result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
			if result.Error != nil || result.RowsAffected == 0 {
				errData := structs.Map(sopsProxyTestResult)
				logger.LogError("ERROR", fileName, result.Error, "Failed to update or create accountProxyStats in DB", errData)
			}

			/*

				SEND SLACK MESSAGE

			*/
			// Create Stats Blob
			statsString := "*URL:* " + fmt.Sprintf("%v", proxyTestSetup["test_url"]) + "\n"
			statsString = statsString + "*Test Mode:* " + fmt.Sprintf("%v", proxyTestSetup["test_mode"]) + "\n"
			statsString = statsString + "*Num Requests:* " + fmt.Sprintf("%v", proxyTestSetup["num_test_requests"]) + "\n\n"
			statsString = statsString + "```Proxy | Success | Valid | Latency |\n"
			for _, workingProxy := range sopsProxyTestResponse.WorkingProxies {
				statsString = statsString + workingProxy.Proxy + " | "
				statsString = statsString + fmt.Sprintf("%v", int(workingProxy.SuccessRate)) + "% | "
				statsString = statsString + fmt.Sprintf("%v", int(workingProxy.ValidRate)) + "% | "
				statsString = statsString + fmt.Sprintf("%v", (math.Round(workingProxy.AvgSuccessLatency*100)/100)) + "s |\n"

			}
			statsString = statsString + "```"

			// Send Slack Message
			headline := "Proxy Test: " + fmt.Sprintf("%v", proxyTestSetup["test_domain"])
			slack.SlackStatsAlert("#proxy-optimizer-tests", headline, statsString)

		} else {

			// User Tests - Get Existing Results and Ammend New Results To Them
			var dbSopsProxyTest models.SopsProxyTestResult
			dbSopsProxyTestResult := db.Where("id = ?", proxyTestSetup["test_id"]).First(&dbSopsProxyTest)
			if dbSopsProxyTestResult.Error == nil {

				var existingSopsProxyTestResults SopsProxyTestResponse
				json.Unmarshal([]byte(dbSopsProxyTest.TestResults), &existingSopsProxyTestResults)

				existingSopsProxyTestResults.TestResults = sopsProxyTestResponse.TestResults
				existingSopsProxyTestResults.WorkingProxies = sopsProxyTestResponse.WorkingProxies
				existingSopsProxyTestResults.CategoryProxies = sopsProxyTestResponse.CategoryProxies

				// Convert the struct to JSON
				jsonSopsProxyTestResponse, err := json.Marshal(existingSopsProxyTestResults)
				if err != nil {
					logger.LogError("ERROR", fileName, err, "Error converting sopsProxyTestResponse to JSON", emptyErrMap)
				}
				sopsProxyTestResultMap["test_results"] = jsonSopsProxyTestResponse

				// Update DB for Successful Test
				result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
				if result.Error != nil || result.RowsAffected == 0 {
					errData := structs.Map(sopsProxyTestResult)
					logger.LogError("ERROR", fileName, result.Error, "Failed to update or create accountProxyStats in DB", errData)
				}

			}

		}

	}

	return
}

func MakeUserProxyTesterRequest(wg *sync.WaitGroup, proxyTestSetup map[string]interface{}, fileName string) {

	defer wg.Done()

	emptyErrMap := make(map[string]interface{})

	// load DB
	var db = db.GetDB()

	logger.LogTextValue("proxyTestSetup", proxyTestSetup)

	client := &http.Client{
		Timeout: 130 * time.Second,
	}

	// Create Request
	var req *http.Request
	var proxyTestEndpoint string

	// proxyTestSetup["api_key"] = "85bb39cd-c5a6-44cc-9221-401e860e52b1"
	postBody, _ := json.Marshal(proxyTestSetup)
	postBodyBytes := bytes.NewBuffer(postBody)

	// Create Proxy Provider Test Request
	proxyTestEndpoint = "https://backend.scrapeops.io/test-user-proxy-settings/v1/"
	req, _ = http.NewRequest("POST", proxyTestEndpoint, postBodyBytes)
	req.Header.Set("Content-Type", "application/json")

	logger.LogTextValue("User Test proxyTestEndpoint", proxyTestEndpoint)
	logger.LogTextValue("User Test Sent Request", proxyTestSetup["test_id"])

	var sopsProxyTestResult models.SopsProxyTestResult

	// Update DB With Test Status
	sopsProxyTestResultMap := map[string]interface{}{}
	sopsProxyTestResultMap["test_status"] = "processing"
	result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
	if result.Error != nil || result.RowsAffected == 0 {
		errData := structs.Map(sopsProxyTestResult)
		logger.LogError("ERROR", fileName, result.Error, "Failed to update sopsProxyTestResultMap in DB", errData)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.LogTextValue("User Test error", err)
		logger.LogTextValue("User Test proxyTestSetup", proxyTestSetup)

		// Update DB for Failed Test
		sopsProxyTestResultMap["test_status"] = "failed"
		result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
		if result.Error != nil || result.RowsAffected == 0 {
			errData := structs.Map(sopsProxyTestResult)
			logger.LogError("ERROR", fileName, result.Error, "Failed to update or create sopsProxyTestResultMap in DB", errData)
		}

	} else {
		logger.LogTextValue("User Test Response", "")
		logger.LogTextValueSpace("User Test response code", resp.StatusCode)

		defer resp.Body.Close()

		// Convert Body To JSON
		var sopsProxyTestResponse SopsProxyTestResponse
		json.NewDecoder(resp.Body).Decode(&sopsProxyTestResponse)

		// logger.LogTextValue("", "")
		logger.LogTextValue("User Test  sopsProxyTestResponse", sopsProxyTestResponse)
		// logger.LogTextValue("sopsProxyTestResultMap", sopsProxyTestResultMap)

		// User Tests - Get Existing Results and Ammend New Results To Them
		var dbSopsProxyTest models.SopsProxyTestResult
		dbSopsProxyTestResult := db.Where("id = ?", proxyTestSetup["test_id"]).First(&dbSopsProxyTest)
		if dbSopsProxyTestResult.Error == nil {

			var existingSopsProxyTestResults SopsProxyTestResponse
			json.Unmarshal([]byte(dbSopsProxyTest.TestResults), &existingSopsProxyTestResults)

			existingSopsProxyTestResults.UserProxyTestResults = sopsProxyTestResponse.UserProxyTestResults

			// Convert the struct to JSON
			jsonSopsProxyTestResponse, err := json.Marshal(existingSopsProxyTestResults)
			if err != nil {
				logger.LogError("ERROR", fileName, err, "Error converting sopsProxyTestResponse to JSON", emptyErrMap)
			}
			sopsProxyTestResultMap["test_results"] = jsonSopsProxyTestResponse

			// Update DB for Successful Test
			result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
			if result.Error != nil || result.RowsAffected == 0 {
				errData := structs.Map(sopsProxyTestResult)
				logger.LogError("ERROR", fileName, result.Error, "Failed to update or create accountProxyStats in DB", errData)
			}

		}

	}

	return
}
