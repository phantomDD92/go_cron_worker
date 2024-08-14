package workers

import (
	"go_proxy_worker/db"
	"go_proxy_worker/logger"
	"go_proxy_worker/slack"
	"go_proxy_worker/utils"
	"os"

	"github.com/go-co-op/gocron"

	// "strconv"
	"strings"
	// "math"
	"fmt"
	"log"
	"time"

	"bytes"
	"encoding/json"
	"net/http"
)

type DDAttributesL2 struct {
	Method           string  `json:"method"`
	Domain           string  `json:"domain"`
	Status           uint    `json:"status"`
	Latency          float64 `json:"latency"`
	FinalProxy       string  `json:"final_proxy"`
	ProxyNumRequests uint    `json:"proxy_num_requests"`
	SopsApiCredits   uint    `json:"sops_api_credits"`
	Country          string  `json:"country"`
	Residential      bool    `json:"residential"`
	Render           bool    `json:"render_js"`
}

type DDAttributesL1 struct {
	AttributesL2 DDAttributesL2 `json:"attributes"`
}

type DDResponseData struct {
	AttributesL1 DDAttributesL1 `json:"attributes"`
}

type DataDogResponse struct {
	Data []DDResponseData `json:"data"`
}

func CronCheckProxyProviderDown() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Minutes().Do(CheckProxyProviderDown)
	s.StartBlocking()
}

func CheckProxyProviderDown() {

	fileName := "cron_check_proxy_provider_down.go"

	// emptyErrMap := make(map[string]interface{})

	var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	// List of Proxy Providers
	proxyProviderList := []string{
		"scraperapi",
		"scrapingant",
		"scrapedo",
		"scrapingfish",
		"scrapingdog",
		"scrapfly",
		"scrapeowl",
		"scrapingbee",
		"zyte_SPM",
		"zyte_SB",
		"infatica_api",
		"brightdata_unlocker",
		"zenrows",
	}

	for _, proxyProvider := range proxyProviderList {

		// Calling Sleep method - getting 429s from
		time.Sleep(4 * time.Second)

		client := &http.Client{
			Timeout: 130 * time.Second,
		}

		// Create Request
		var req *http.Request

		/*

				POST BODY - Example:

				post_body = {
				  "filter": {
					"from": "now-3d",
					"to": "now",
					"query": "service:scrapeops-go-proxy-std-out @log_type:user_log @email:send2kust@gmail.com -status:warn"
				  },
				"page": {
					"limit":1000
				  },
			}

		*/

		queryString := "service:scrapeops-go-proxy-std-out @log_type:user_log -status:warn @final_proxy:" + proxyProvider
		filterMap := map[string]interface{}{"from": "now-30m", "to": "now", "query": queryString}
		pageMap := map[string]interface{}{"limit": 1000}
		postBodyMap := map[string]interface{}{"filter": filterMap, "page": pageMap}
		postBody, _ := json.Marshal(postBodyMap)
		postBodyBytes := bytes.NewBuffer(postBody)

		// Create Request
		datadogEndpoint := "https://api.datadoghq.eu/api/v2/logs/events/search"
		req, _ = http.NewRequest("POST", datadogEndpoint, postBodyBytes)

		// Add Headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("DD-API-KEY", os.Getenv("DD_API_KEY"))
		req.Header.Set("DD-APPLICATION-KEY", os.Getenv("DD_APPLICATION_KEY"))
		// req.SetBasicAuth(apiKey, "")

		// Make Request
		resp, err := client.Do(req)
		if err != nil {
			log.Println("error", queryString)
		}

		// else {
		// 	log.Println("dd response code", resp.StatusCode)
		// }

		// Parse Result
		if err == nil && resp.StatusCode == 200 {

			defer resp.Body.Close()

			// Convert Body To JSON
			var dataDogResponse DataDogResponse
			json.NewDecoder(resp.Body).Decode(&dataDogResponse)

			// log.Println("dataDogResponse", dataDogResponse)

			var datadogResults []DDAttributesL2
			for _, row := range dataDogResponse.Data {
				datadogResults = append(datadogResults, row.AttributesL1.AttributesL2)
			}

			// log.Println("proxyProvider", proxyProvider)
			// log.Println("datadogResults", len(datadogResults))

			if len(datadogResults) > 0 {

				// Count Successful Requests
				var successful uint
				var uniqueDomains []string
				for _, request := range datadogResults {
					if request.Status == 200 || request.Status == 404 {
						successful = successful + 1
						if utils.StringInSlice(request.Domain, uniqueDomains) == false {
							uniqueDomains = append(uniqueDomains, request.Domain)
						}
					}
				}

				// log.Println("successful", successful)
				// log.Println("uniqueDomains", len(uniqueDomains))

				// All Requests Are Failing
				if len(datadogResults) > 900 && successful == 0 && len(uniqueDomains) > 3 {
					log.Println("all failed: ", proxyProvider)

					// Message Slack
					headline := strings.Title(proxyProvider) + ": Down"
					msg := "100% failure rate over " + fmt.Sprintf("%v", len(uniqueDomains)) + " domains"
					slack.SlackStatsAlert("#proxy-provider-down", headline, msg)

					/*

						REMOVING PROXY PROVIDER
						Uses 2 keys:
						 - proxyProviderDownKey
						 - proxyProviderDownRetestKey

						 If the proxyProviderDownKey isn't set to true, then it will check the proxyProviderDownRetestKey
						 which is designed to give a 5 minute window for the system to start resending requests to the
						 proxy. So that it can gain more data and check if the proxy is still down.

						 Otherwise, no new requests would be sent to the proxy after removing it and the system would
						 still think that the proxy provider is down.

					*/

					// Check Redis For Blocked Proxy Already
					proxyProviderDownKey := "proxyProviderDown?proxy=" + proxyProvider
					redisProxyDownResult, err1 := coreProxyRedisClient.Get(redisContext, proxyProviderDownKey).Result()
					if redisProxyDownResult == "" && err1 == nil {

						proxyProviderDownRetestKey := "proxyProviderDown?proxy=" + proxyProvider + "&retest=true"
						redisProxyDownRetestResult, err2 := coreProxyRedisClient.Get(redisContext, proxyProviderDownRetestKey).Result()
						if redisProxyDownRetestResult == "" && err2 == nil {

							// Tell System That Proxy Provider Is Down - Done to tell proxy manager to find which one is down
							err = coreProxyRedisClient.Set(redisContext, "proxyProviderDown?overall=true", "true", 15*60*time.Second).Err()
							if err != nil {
								logger.LogError("ERROR", fileName, err, "failed to deactivate proxy provider in Redis", map[string]interface{}{
									"proxyProviderDownKey": proxyProviderDownKey,
								})
							}

							// Pull Proxy Out Of Rotation
							err = coreProxyRedisClient.Set(redisContext, proxyProviderDownKey, "true", 15*60*time.Second).Err()
							if err != nil {
								logger.LogError("ERROR", fileName, err, "failed to deactivate proxy provider in Redis", map[string]interface{}{
									"proxyProviderDownKey": proxyProviderDownKey,
								})
							}

							// Set Retest Limit
							err = coreProxyRedisClient.Set(redisContext, proxyProviderDownRetestKey, "true", 20*60*time.Second).Err()
							if err != nil {
								logger.LogError("ERROR", fileName, err, "failed to set proxy provider retest in Redis", map[string]interface{}{
									"proxyProviderDownKey":       proxyProviderDownKey,
									"proxyProviderDownRetestKey": proxyProviderDownRetestKey,
								})
							}

							// Message Slack - Pull Proxy Out of Rotation
							headline := strings.Title(proxyProvider) + ": Removed From Rotation"
							msg := strings.Title(proxyProvider) + " removed from proxy rotation for 15 minutes."
							slack.SlackStatsAlert("#proxy-provider-down", headline, msg)

						} else {

							// Message Slack - Pull Proxy Out of Rotation
							headline := strings.Title(proxyProvider) + ": Retesting"
							msg := strings.Title(proxyProvider) + " enabled for 5 minutes of retesting."
							slack.SlackStatsAlert("#proxy-provider-down", headline, msg)
						}

					}

				}

			}

		}

	}

}
