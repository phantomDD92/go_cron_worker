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
	TestResults     []map[string]interface{} `json:"test_results"`
	WorkingProxies  []WorkingProxy           `json:"working_proxies"`
	CategoryProxies []map[string]interface{} `json:"category_proxies"`
}

type WorkingProxy struct {
	Proxy             string            `json:"proxy"`
	ConcurrencyLimit  int               `json:"concurrency_limit"`
	Type              string            `json:"type"`
	Features          map[string]string `json:"features"`
	TotalRequests     int               `json:"total_requests"`
	Successful        int               `json:"successful"`
	Test1Passed       int               `json:"test_1_passed"`
	TotalLatency      int               `json:"total_latency"`
	SuccessRate       float64           `json:"success_rate"`
	ValidRate         float64           `json:"valid_rate"`
	AvgSuccessLatency float64           `json:"avg_success_latency"`
	CPM               int               `json:"CPM"`
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

	// proxyTestSetup["api_key"] = "85bb39cd-c5a6-44cc-9221-401e860e52b1"
	postBody, _ := json.Marshal(proxyTestSetup)
	postBodyBytes := bytes.NewBuffer(postBody)

	if proxyTestSetup["test_type"] == "test_internal_proxy_pools" {
		// Create Proxy Provider Test Request
		proxyTestEndpoint := "https://backend.scrapeops.io/test-proxies-providers/v2/"
		req, _ = http.NewRequest("POST", proxyTestEndpoint, postBodyBytes)
		req.Header.Set("Content-Type", "application/json")
	} else {
		// test_user_proxy_settings
		// Create User Proxy Settings Test Request
		proxyTestEndpoint := "https://backend.scrapeops.io/test-user-proxy-settings/v1/"
		req, _ = http.NewRequest("POST", proxyTestEndpoint, postBodyBytes)
		req.Header.Set("Content-Type", "application/json")
	}

	// Update DB With Test Status
	sopsProxyTestResultMap := map[string]interface{}{
		"test_status": "processing",
	}

	var sopsProxyTestResult models.SopsProxyTestResult
	result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
	if result.Error != nil || result.RowsAffected == 0 {
		errData := structs.Map(sopsProxyTestResult)
		logger.LogError("ERROR", fileName, result.Error, "Failed to update or create accountProxyStats in DB", errData)
	}

	logger.LogTextValue("Sent Request", proxyTestSetup["test_id"])

	resp, err := client.Do(req)
	if err != nil {
		logger.LogTextValue("error", proxyTestSetup)

		// Update DB for Failed Test
		sopsProxyTestResultMap["status"] = "failed"
		result := db.Model(&sopsProxyTestResult).Where("id = ?", proxyTestSetup["test_id"]).Updates(sopsProxyTestResultMap)
		if result.Error != nil || result.RowsAffected == 0 {
			errData := structs.Map(sopsProxyTestResult)
			logger.LogError("ERROR", fileName, result.Error, "Failed to update or create accountProxyStats in DB", errData)
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

		// Update DB for Failed Test
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

	}

	return
}
