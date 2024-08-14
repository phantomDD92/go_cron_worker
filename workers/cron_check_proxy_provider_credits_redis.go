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

func CronCheckProxyProviderCreditsRedis() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Seconds().Do(CheckProxyProviderCreditsRedis)
	s.StartBlocking()
}

// type ScrapingAntAccount struct {
// 	PlanName					string `json:"plan_name"`
// 	PlanTotalCredits			int `json:"plan_total_credits"`
// 	RemainedCredits				int `json:"remained_credits"`
// }

func CheckProxyProviderCreditsRedis() {

	fileName := "cron_check_proxy_provider_credits.go"

	// emptyErrMap := make(map[string]interface{})

	// Redis Details
	var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	// var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	/*

		SCRAPINGANT

	*/

	scrapingantAccountEndpoint := "https://api.scrapingant.com/v1/usage?x-api-key="

	scrapingantAPIKeys := []string{os.Getenv("SCRAPINGANT_INFO_API_KEY"), os.Getenv("SCRAPINGANT_IAN_API_KEY")}

	for _, apiKey := range scrapingantAPIKeys {

		url := scrapingantAccountEndpoint + apiKey
		req, _ := http.NewRequest("GET", url, nil)

		client := &http.Client{
			Timeout: 130 * time.Second,
		}

		// Make Request
		resp, err := client.Do(req)

		if err == nil && resp.StatusCode == 200 {

			defer resp.Body.Close()

			// Convert Body To JSON
			var scrapingAntAccount ScrapingAntAccount
			json.NewDecoder(resp.Body).Decode(&scrapingAntAccount)

			proxyProviderApiKeyActiveKey := "proxyProviderApiKeyActive?proxyProvider=scrapingant&api_key=" + apiKey

			if scrapingAntAccount.RemainedCredits == 0 {

				// Update Redis To Deactivate API Key
				err = coreProxyRedisClient.Set(redisContext, proxyProviderApiKeyActiveKey, "false", 60*60*time.Second).Err()
				if err != nil {
					logger.LogError("ERROR", fileName, err, "failed to deactivate proxy provider in Redis", map[string]interface{}{
						"proxyProviderApiKeyActiveKey": proxyProviderApiKeyActiveKey,
					})
				}

			} else {

				// Update Redis To Activate API Key
				err = coreProxyRedisClient.Set(redisContext, proxyProviderApiKeyActiveKey, "true", 60*60*time.Second).Err()
				if err != nil {
					logger.LogError("ERROR", fileName, err, "failed to deactivate proxy provider in Redis", map[string]interface{}{
						"proxyProviderApiKeyActiveKey": proxyProviderApiKeyActiveKey,
					})
				}

			}

		}
	}

}
