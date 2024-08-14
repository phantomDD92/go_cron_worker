package workers

import (
	"bytes"
	"encoding/json"
	"go_proxy_worker/logger"
	"go_proxy_worker/slack"
	"os"
	"time"

	"github.com/go-co-op/gocron"

	// "log"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

func CronCheckProxyApiProfitability() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(24).Hours().Do(CheckProxyApiProfitability)
	s.StartBlocking()
}

// Define structs to match the JSON response structure
type RawResponse struct {
	Data []struct {
		Attributes struct {
			Attributes struct {
				Method                  string  `json:"method"`
				Domain                  string  `json:"domain"`
				Status                  int     `json:"status"`
				Duration                float64 `json:"duration"`
				FinalProxy              string  `json:"final_proxy"`
				ProxyNumRequests        int     `json:"proxy_num_requests"`
				SopsAPICredits          int     `json:"sops_api_credits"`
				Country                 string  `json:"country"`
				Residential             bool    `json:"residential"`
				RenderJS                bool    `json:"render_js"`
				RequestProfit           float64 `json:"sops_request_profit"`
				RequestProfitPercentage float64 `json:"sops_request_profit_%"`
			} `json:"attributes"`
			Timestamp string `json:"timestamp"`
		} `json:"attributes"`
	} `json:"data"`
}

// Define a struct for the condensed data
type CondensedData struct {
	// Method           string  `json:"method"`
	// Domain           string  `json:"domain"`
	// Status           int  	 `json:"status"`
	// Latency          float64 `json:"latency"`
	// Timestamp        string  `json:"timestamp"`
	// FinalProxy       string  `json:"final_proxy"`
	// ProxyNumRequests int     `json:"proxy_num_requests"`
	// SopsAPICredits   int     `json:"sops_api_credits"`
	// Country          string  `json:"country"`
	// Residential      bool    `json:"residential"`
	// RenderJS         bool    `json:"render_js"`
	RequestProfit           float64 `json:"sops_request_profit"`
	RequestProfitPercentage float64 `json:"sops_request_profit_%"`
}

func CheckProxyApiProfitability() {

	fileName := "cron_check_proxy_api_profitability.go"

	emptyErrMap := make(map[string]interface{})

	// Date Info
	now := time.Now()
	dayStartTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayStartDateString := dayStartTime.Format("2006-01-02")

	// Define headers
	headers := map[string]string{
		"Content-Type":       "application/json",
		"DD-API-KEY":         os.Getenv("DD_API_KEY"),
		"DD-APPLICATION-KEY": os.Getenv("DD_APPLICATION_KEY"),
	}

	// Define the 24-hour period (24 hours ago from now until now)
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	// Step through every 5 minutes
	step := 5 * time.Minute

	var condensedData []CondensedData

	for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.Add(step) {
		fromTime := currentTime
		toTime := currentTime.Add(step)

		// Format "from" and "to" as strings for the post body
		fromTimeStr := fromTime.Format(time.RFC3339)
		toTimeStr := toTime.Format(time.RFC3339)

		postBody := map[string]interface{}{
			"filter": map[string]string{
				"from":  fromTimeStr,
				"to":    toTimeStr,
				"query": "service:scrapeops-go-proxy-std-out -@proxy_plan_id:3 -status:(warn or error)", // Replace with your actual query
			},
			"page": map[string]int{
				"limit": 1000, // Adjust as needed
			},
		}

		// Here, replace this with the actual function call to send your request
		logger.LogTextValue("Sending request with From: ", fromTimeStr)
		// For example: sendRequest(postBody)

		// Encode the post body to JSON
		bodyBytes, err := json.Marshal(postBody)
		if err != nil {
			logger.LogError("ERROR", fileName, err, "Error encoding JSON", emptyErrMap)
			// return
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "https://api.datadoghq.eu/api/v2/logs/events/search", bytes.NewBuffer(bodyBytes))
		if err != nil {
			logger.LogError("ERROR", fileName, err, "Error creating request", emptyErrMap)
			// return
		}

		// Add headers to the request
		for key, value := range headers {
			req.Header.Add(key, value)
		}

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.LogError("ERROR", fileName, err, "Error sending request", emptyErrMap)
			// return
		}

		logger.LogTextValue("Response status:", resp.Status)

		if err == nil && resp.StatusCode == 200 {

			defer resp.Body.Close()

			// Read the response body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.LogError("ERROR", fileName, err, "Error reading response body", emptyErrMap)
			}

			// Unmarshal the JSON response into RawResponse struct
			var rawResponse RawResponse
			err = json.Unmarshal(body, &rawResponse)
			if err != nil {
				logger.LogError("ERROR", fileName, err, "Error unmarshaling JSON", emptyErrMap)
			}

			// Process the data and extract condensed data
			for _, row := range rawResponse.Data {
				attrL1 := row.Attributes
				attrL2 := attrL1.Attributes
				condensedData = append(condensedData, CondensedData{
					// Method:           attrL2.Method,
					// Domain:           attrL2.Domain,
					// Status:           attrL2.Status,
					// Latency:          attrL2.Duration,
					// Timestamp:        attrL1.Timestamp,
					// FinalProxy:       attrL2.FinalProxy,
					// ProxyNumRequests: attrL2.ProxyNumRequests,
					// SopsAPICredits:   attrL2.SopsAPICredits,
					// Country:          attrL2.Country,
					// Residential:      attrL2.Residential,
					// RenderJS:         attrL2.RenderJS,
					RequestProfit:           attrL2.RequestProfit,
					RequestProfitPercentage: attrL2.RequestProfitPercentage,
				})
			}

			logger.LogTextValue("num_raw_logs:", len(rawResponse.Data))
			logger.LogTextValue("num_logs:", len(condensedData))

		}

		time.Sleep(2 * time.Second) // Wait for 3 seconds

	}

	// Output the number of logs and optionally the condensed data
	logger.LogTextValue("num_logs:", len(condensedData))

	if len(condensedData) > 0 {

		var totalProfit float64
		var totalProfitPercentage float64
		var unprofitableRequests float64

		for _, data := range condensedData {
			totalProfit += data.RequestProfit
			totalProfitPercentage += data.RequestProfitPercentage
			if data.RequestProfit < 0 {
				unprofitableRequests = unprofitableRequests + 1
			}
		}

		totalRequestsFloat64 := float64(len(condensedData))
		meanProfit := totalProfit / totalRequestsFloat64
		meanProfitPercentage := totalProfitPercentage / totalRequestsFloat64
		unprofitableRequestPercentage := (unprofitableRequests / totalRequestsFloat64) * 100

		logger.LogTextValue("meanProfit CPM: $", meanProfit*1000000)
		logger.LogTextValue("meanProfitPercentage:", meanProfitPercentage)
		logger.LogTextValue("unprofitableRequestPercentage:", unprofitableRequestPercentage)

		// Create Stats Blob
		statsString := "```Mean Profit CPM: $" + fmt.Sprintf("%v", math.Round(meanProfit*1000000*100)/100) + "/million \n"
		statsString = statsString + "Mean Profit Percentage: " + fmt.Sprintf("%v", math.Round(meanProfitPercentage*100)/100) + "% \n"
		statsString = statsString + "Unprofitable Requests Percentage: " + fmt.Sprintf("%v", math.Round(unprofitableRequestPercentage*100)/100) + "% \n"
		statsString = statsString + "Number of Logs: " + fmt.Sprintf("%v", len(condensedData)) + " \n"
		statsString = statsString + "```"

		// Send Slack Message
		headline := "Proxy API Profitability: " + dayStartDateString
		slack.SlackStatsAlert("#proxy-profitability", headline, statsString)

	}

}
