package workers

import (
	"go_proxy_worker/logger"
	"go_proxy_worker/utils"
	"os"

	"github.com/go-co-op/gocron"

	// "go_proxy_worker/dbRedisQueries"

	"encoding/json"
	// "io/ioutil"

	"go_proxy_worker/db"
	"net/http"
	"time"
	// "log"
)

func CronSyncProxyConcurrency() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(15).Seconds().Do(SyncProxyConcurrency)
	s.StartBlocking()
}

// type ScrapingbeeAccount struct {
// 	MaxApiCredit				int `json:"max_api_credit"`
// 	UsedApiCredit				int `json:"used_api_credit"`
// 	MaxConcurrency				int `json:"max_concurrency"`
// 	CurrentConcurrency			int `json:"current_concurrency"`
// }

// type ScraperapiAccount struct {
// 	RequestLimit				int `json:"requestLimit"`
// 	RequestCount				int `json:"requestCount"`
// 	ConcurrencyLimit			int `json:"concurrencyLimit"`
// 	ConcurrentRequests			int `json:"concurrentRequests"`
// }

// type ScrapingdogAccount struct {
// 	RequestLimit				int `json:"requestLimit"`
// 	RequestUsed					int `json:"requestUsed"`
// 	ActiveThread				int `json:"activeThread"`
// }

// type ScrapeowlAccount struct {
// 	Credits						int `json:"credits"`
// 	CreditsUsed					int `json:"credits_used"`
// 	Requests					int `json:"requests"`
// 	FailedRequests				int `json:"failed_requests"`
// 	ConcurrencyLimit			int `json:"concurrency_limit"`
// 	ConcurrentRequests			int `json:"concurrent_requests"`
// }

// type ScrapflyAccount struct {
// 	Subscription		ScrapflySubscription `json:"subscription"`
// }

// type ScrapflySubscription struct {
// 	Usage			ScrapflyUsage `json:"usage"`
// }

// type ScrapflyUsage struct {
// 	Scrape			ScrapflyScrape `json:"scrape"`
// }

// type ScrapflyScrape struct {
// 	ConcurrentUsage			int `json:"concurrent_usage"`
// 	ConcurrentLimit			int `json:"concurrent_limit"`
// 	ConcurrentRemaining			int `json:"concurrent_remaining"`
// }

func SyncProxyConcurrency() {

	fileName := "cron_sync_proxy_concurrency.go"

	// emptyErrMap := make(map[string]interface{})

	// Redis Details
	var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	// var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	/*

		SCRAPINGBEE

	*/

	scrapingbeeApiKey := os.Getenv("SCRAPINGBEE_API_KEY")
	scrapingbeeAccountEndpoint := "https://app.scrapingbee.com/api/v1/usage?api_key=" + scrapingbeeApiKey

	req, _ := http.NewRequest("GET", scrapingbeeAccountEndpoint, nil)

	client := &http.Client{
		Timeout: 130 * time.Second,
	}

	// Make Request
	resp, err := client.Do(req)

	if err == nil && resp.StatusCode == 200 {

		defer resp.Body.Close()

		// body, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	logger.LogError("ERROR", fileName, err, "Error reading response body", emptyErrMap)
		// }

		// Convert Body To JSON
		var scrapingbeeAccount ScrapingbeeAccount
		json.NewDecoder(resp.Body).Decode(&scrapingbeeAccount)

		// log.Println("scrapingbeeAccount.CurrentConcurrency", scrapingbeeAccount.CurrentConcurrency)

		// Update Redis
		redisProxyConcurrencyKey := "proxyConcurrency?proxyProvider=scrapingbee"

		err = coreProxyRedisClient.Set(redisContext, redisProxyConcurrencyKey, scrapingbeeAccount.CurrentConcurrency, 5*60*time.Second).Err()
		if err != nil {
			logger.LogError("ERROR", fileName, err, "failed to sync proxy concurrency in Redis", map[string]interface{}{
				"proxy":              "scrapingbee",
				"CurrentConcurrency": scrapingbeeAccount.CurrentConcurrency,
			})
		}
	}

	/*

		SCRAPERAPI

	*/

	scraperapiApiKey := os.Getenv("SCRAPERAPI_API_KEY")
	scraperapiAccountEndpoint := "http://api.scraperapi.com/account?api_key=" + scraperapiApiKey

	req, _ = http.NewRequest("GET", scraperapiAccountEndpoint, nil)

	client = &http.Client{
		Timeout: 130 * time.Second,
	}

	// Make Request
	resp, err = client.Do(req)

	if err == nil && resp.StatusCode == 200 {

		defer resp.Body.Close()

		// Convert Body To JSON
		var scraperapiAccount ScraperapiAccount
		json.NewDecoder(resp.Body).Decode(&scraperapiAccount)

		// Update Redis
		redisProxyConcurrencyKey := "proxyConcurrency?proxyProvider=scraperapi"

		err = coreProxyRedisClient.Set(redisContext, redisProxyConcurrencyKey, scraperapiAccount.ConcurrentRequests, 5*60*time.Second).Err()
		if err != nil {
			logger.LogError("ERROR", fileName, err, "failed to sync proxy concurrency in Redis", map[string]interface{}{
				"proxy":              "scraperapi",
				"ConcurrentRequests": scraperapiAccount.ConcurrentRequests,
			})
		}
	}

	/*

		SCRAPINGDOG

	*/

	scrapingdogApiKey := os.Getenv("SCRAPINGDOG_API_KEY")
	scrapingdogAccountEndpoint := "https://api.scrapingdog.com/account?api_key=" + scrapingdogApiKey

	req, _ = http.NewRequest("GET", scrapingdogAccountEndpoint, nil)

	client = &http.Client{
		Timeout: 130 * time.Second,
	}

	// Make Request
	resp, err = client.Do(req)

	if err == nil && resp.StatusCode == 200 {

		defer resp.Body.Close()

		// Convert Body To JSON
		var scrapingdogAccount ScrapingdogAccount
		json.NewDecoder(resp.Body).Decode(&scrapingdogAccount)

		// Update Redis
		redisProxyConcurrencyKey := "proxyConcurrency?proxyProvider=scrapingdog"

		err = coreProxyRedisClient.Set(redisContext, redisProxyConcurrencyKey, scrapingdogAccount.ActiveThread, 5*60*time.Second).Err()
		if err != nil {
			logger.LogError("ERROR", fileName, err, "failed to sync proxy concurrency in Redis", map[string]interface{}{
				"proxy":              "scrapingdog",
				"ConcurrentRequests": scrapingdogAccount.ActiveThread,
			})
		}
	}

	/*

		SCRAPFLY

	*/

	scrapflyApiKey := os.Getenv("SCRAPFLY_API_KEY")
	scrapflyAccountEndpoint := "https://api.scrapfly.io/account?key=" + scrapflyApiKey

	req, _ = http.NewRequest("GET", scrapflyAccountEndpoint, nil)

	client = &http.Client{
		Timeout: 130 * time.Second,
	}

	// Make Request
	resp, err = client.Do(req)

	if err == nil && resp.StatusCode == 200 {

		defer resp.Body.Close()

		// Convert Body To JSON
		var scrapflyAccount ScrapflyAccount
		json.NewDecoder(resp.Body).Decode(&scrapflyAccount)

		// log.Println("Scrapfly Scrape", scrapflyAccount.Subscription.Usage.Scrape)
		// log.Println("Scrapfly Concurrency", scrapflyAccount.Subscription.Usage.Scrape.ConcurrentUsage)

		// Update Redis
		redisProxyConcurrencyKey := "proxyConcurrency?proxyProvider=scrapfly"

		err = coreProxyRedisClient.Set(redisContext, redisProxyConcurrencyKey, scrapflyAccount.Subscription.Usage.Scrape.ConcurrentUsage, 5*60*time.Second).Err()
		if err != nil {
			logger.LogError("ERROR", fileName, err, "failed to sync proxy concurrency in Redis", map[string]interface{}{
				"proxy":              "scrapfly",
				"ConcurrentRequests": scrapflyAccount.Subscription.Usage.Scrape.ConcurrentUsage,
			})
		}
	}

	/*

		SCRAPEOWL

	*/

	scrapeOwlApiKeyList := []string{os.Getenv("SCRAPEOWL_API_KEY")}
	for _, scrapeowlApiKey := range scrapeOwlApiKeyList {

		scrapeAccountEndpoint := "https://api.scrapeowl.com/v1/usage?api_key=" + scrapeowlApiKey

		req, _ = http.NewRequest("GET", scrapeAccountEndpoint, nil)

		client = &http.Client{
			Timeout: 130 * time.Second,
		}

		// Make Request
		resp, err = client.Do(req)

		if err == nil && resp.StatusCode == 200 {

			defer resp.Body.Close()

			// Convert Body To JSON
			var scrapeowlAccount ScrapeowlAccount
			json.NewDecoder(resp.Body).Decode(&scrapeowlAccount)

			// log.Println("Scrapingowl Scrape", scrapeowlAccount.CreditsUsed)
			// log.Println("Scrapingowl Concurrency", scrapeowlAccount.ConcurrentRequests)

			// Update Redis
			redisProxyConcurrencyKey := "proxyConcurrency?proxyProvider=scrapeowl&apiKey=" + scrapeowlApiKey

			err = coreProxyRedisClient.Set(redisContext, redisProxyConcurrencyKey, scrapeowlAccount.ConcurrentRequests, 5*60*time.Second).Err()
			if err != nil {
				logger.LogError("ERROR", fileName, err, "failed to sync proxy concurrency in Redis", map[string]interface{}{
					"proxy":              "scrapeowl__" + scrapeowlApiKey,
					"ConcurrentRequests": scrapeowlAccount.ConcurrentRequests,
				})
			}
		}

	}

}
