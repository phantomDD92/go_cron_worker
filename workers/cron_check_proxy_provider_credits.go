package workers

import (
	"fmt"
	"go_proxy_worker/slack"
	"math"
	"os"

	"github.com/go-co-op/gocron"
	// "go_proxy_worker/logger"
	// "go_proxy_worker/utils"

	"encoding/json"
	// "io/ioutil"

	// "go_proxy_worker/db"
	"log"
	"net/http"
	"sort"
	"time"
)

func CronCheckProxyProviderCredits() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(12).Hours().Do(CheckProxyProviderCredits)
	s.StartBlocking()
}

type ScrapedoAccount struct {
	MaxMonthlyRequest          int `json:"MaxMonthlyRequest"`
	RemainingMonthlyRequest    int `json:"RemainingMonthlyRequest"`
	ConcurrentRequest          int `json:"ConcurrentRequest"`
	RemainingConcurrentRequest int `json:"RemainingConcurrentRequest"`
}

type ScrapingbeeAccount struct {
	MaxApiCredit                  int    `json:"max_api_credit"`
	UsedApiCredit                 int    `json:"used_api_credit"`
	MaxConcurrency                int    `json:"max_concurrency"`
	CurrentConcurrency            int    `json:"current_concurrency"`
	RenewalSubscriptionDateString string `json:"renewal_subscription_date"`
	// RenewalSubscriptionDate       time.Time `json:"renewal_subscription_date"`
}

type ScraperapiAccount struct {
	RequestLimit       int       `json:"requestLimit"`
	RequestCount       int       `json:"requestCount"`
	ConcurrencyLimit   int       `json:"concurrencyLimit"`
	ConcurrentRequests int       `json:"concurrentRequests"`
	SubscriptionDate   time.Time `json:"subscriptionDate"`
}

type ScrapingdogAccount struct {
	RequestLimit int `json:"requestLimit"`
	RequestUsed  int `json:"requestUsed"`
	ActiveThread int `json:"activeThread"`
}

type ScrapeowlAccount struct {
	Credits            int `json:"credits"`
	CreditsUsed        int `json:"credits_used"`
	Requests           int `json:"requests"`
	FailedRequests     int `json:"failed_requests"`
	ConcurrencyLimit   int `json:"concurrency_limit"`
	ConcurrentRequests int `json:"concurrent_requests"`
}

type ScrapingAntAccount struct {
	PlanName         string `json:"plan_name"`
	PlanTotalCredits int    `json:"plan_total_credits"`
	RemainedCredits  int    `json:"remained_credits"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
}

// Scrapfly
type ScrapflyScrape struct {
	Current             int `json:"current"`
	Limit               int `json:"limit"`
	Remaining           int `json:"remaining"`
	ConcurrentUsage     int `json:"concurrent_usage"`
	ConcurrentLimit     int `json:"concurrent_limit"`
	ConcurrentRemaining int `json:"concurrent_remaining"`
}

type ScrapflyPeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ScrapflyUsage struct {
	Scrape ScrapflyScrape `json:"scrape"`
}

type ScrapflySubscription struct {
	Usage  ScrapflyUsage  `json:"usage"`
	Period ScrapflyPeriod `json:"period"`
}

type ScrapflyAccount struct {
	Subscription ScrapflySubscription `json:"subscription"`
}

// Zyte
type ZyteResults struct {
	Failed int `json:"failed"`
	Clean  int `json:"clean"`
}

type ZyteResponse struct {
	After   string        `json:"after"`
	Limit   int           `json:"limit"`
	Results []ZyteResults `json:"results"`
}

// Zenrows
type ZenrowsAccount struct {
	ApiCreditUsage int `json:"api_credit_usage"`
	ApiCreditLimit int `json:"api_credit_limit"`
}

type ProxyProviderStats struct {
	Name             string    `json:"name"`
	CreditLimit      int       `json:"credit_limit"`
	UsedCredits      int       `json:"used_credits"`
	UsedPercentage   float64   `json:"used_percentage"`
	RenewalDate      time.Time `json:"renewal_date"`
	DaysUntilRenewal int       `json:"days_until_renewal"`
}

func CheckProxyProviderCredits() {

	// fileName := "cron_check_proxy_provider_credits.go"

	// emptyErrMap := make(map[string]interface{})

	renewalDateLayout := "2006-01-02T15:04:05.000000"

	// Proxy Provider Stats
	var proxyProviderStatsArray []ProxyProviderStats

	now := time.Now()
	dayStartTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayStartDateString := dayStartTime.Format("2006-01-02 15:04")
	// dayStartDate, _ := time.Parse("2006-01-02 15:04", dayStartDateString)

	// timeZoneLocation, err := time.LoadLocation("UTC")
	// if err != nil {
	// 	log.Println(err)
	// }

	// year := now.Year()

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

		// Convert Body To JSON
		var scrapingbeeAccount ScrapingbeeAccount
		json.NewDecoder(resp.Body).Decode(&scrapingbeeAccount)

		proxyStats := ProxyProviderStats{
			Name:           "Scrapingbee",
			CreditLimit:    scrapingbeeAccount.MaxApiCredit,
			UsedCredits:    scrapingbeeAccount.UsedApiCredit,
			UsedPercentage: float64(scrapingbeeAccount.UsedApiCredit) / float64(scrapingbeeAccount.MaxApiCredit),
			// RenewalDate:      scrapingbeeAccount.RenewalSubscriptionDate,
			// DaysUntilRenewal: int(scrapingbeeAccount.RenewalSubscriptionDate.Sub(now).Hours() / 24),
		}

		parsedTime, err := time.Parse(renewalDateLayout, scrapingbeeAccount.RenewalSubscriptionDateString)
		if err != nil {
			log.Println("Error Parsing Scrapingbee Date:", scrapingbeeAccount.RenewalSubscriptionDateString)
		} else {
			proxyStats.RenewalDate = parsedTime
			proxyStats.DaysUntilRenewal = int(parsedTime.Sub(now).Hours() / 24)
		}

		// log.Println("Scrapingbee RenewalSubscriptionDateString:", scrapingbeeAccount.RenewalSubscriptionDateString)

		proxyProviderStatsArray = append(proxyProviderStatsArray, proxyStats)

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

		proxyStats := ProxyProviderStats{
			Name:             "ScraperAPI",
			CreditLimit:      scraperapiAccount.RequestLimit,
			UsedCredits:      scraperapiAccount.RequestCount,
			UsedPercentage:   float64(scraperapiAccount.RequestCount) / float64(scraperapiAccount.RequestLimit),
			RenewalDate:      scraperapiAccount.SubscriptionDate,
			DaysUntilRenewal: 31 + int(scraperapiAccount.SubscriptionDate.Sub(now).Hours()/24),
		}

		// log.Println("scraperapiAccount.SubscriptionDate:", scraperapiAccount.SubscriptionDate)
		// log.Println("scraperapiAccount.SubscriptionDate.Sub(now).Hours()/24:", scraperapiAccount.SubscriptionDate.Sub(now).Hours()/24)

		proxyProviderStatsArray = append(proxyProviderStatsArray, proxyStats)

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

		proxyStats := ProxyProviderStats{
			Name:             "ScrapingDog",
			CreditLimit:      scrapingdogAccount.RequestLimit,
			UsedCredits:      scrapingdogAccount.RequestUsed,
			UsedPercentage:   float64(scrapingdogAccount.RequestUsed) / float64(scrapingdogAccount.RequestLimit),
			DaysUntilRenewal: 999,
		}

		proxyProviderStatsArray = append(proxyProviderStatsArray, proxyStats)
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

		// log.Println("scrapflyAccount", scrapflyAccount)
		// log.Println("Scrapfly Period End", scrapflyAccount.Subscription.Period.End)
		// log.Println("Scrapfly Usage", scrapflyAccount.Subscription.Usage)
		// log.Println("Scrapfly Scrape", scrapflyAccount.Subscription.Usage.Scrape)
		// log.Println("Scrapfly Limit", scrapflyAccount.Subscription.Usage.Scrape.Limit)
		// log.Println("Scrapfly Current", scrapflyAccount.Subscription.Usage.Scrape.Current)

		// Parse the start and end dates
		// startDate, err := time.Parse("2006-01-02 15:04:05", account.Subscription.Period.Start)
		// if err != nil {
		// 	fmt.Println("Error parsing start date:", err)
		// 	return
		// }

		endDate, err := time.Parse("2006-01-02 15:04:05", scrapflyAccount.Subscription.Period.End)
		if err != nil {
			log.Println("Scrapfly Error parsing end date:", err)
		}

		proxyStats := ProxyProviderStats{
			Name:             "Scrapfly",
			CreditLimit:      scrapflyAccount.Subscription.Usage.Scrape.Limit,
			UsedCredits:      scrapflyAccount.Subscription.Usage.Scrape.Current,
			UsedPercentage:   float64(scrapflyAccount.Subscription.Usage.Scrape.Current) / float64(scrapflyAccount.Subscription.Usage.Scrape.Limit),
			RenewalDate:      endDate,
			DaysUntilRenewal: int(endDate.Sub(now).Hours() / 24),
		}

		proxyProviderStatsArray = append(proxyProviderStatsArray, proxyStats)
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

			var subAccount string
			if scrapeowlApiKey == "d18e85180c3f8fcfd41fb9d3d5efda" {
				subAccount = "(info)"
			} else if scrapeowlApiKey == "3844aa9d5c71e5c0034b01f27df879" {
				subAccount = "(ian)"
			}

			proxyStats := ProxyProviderStats{
				Name:             "ScrapeOwl" + subAccount,
				CreditLimit:      scrapeowlAccount.Credits,
				UsedCredits:      scrapeowlAccount.CreditsUsed,
				UsedPercentage:   float64(scrapeowlAccount.CreditsUsed) / float64(scrapeowlAccount.Credits),
				DaysUntilRenewal: 999,
			}

			proxyProviderStatsArray = append(proxyProviderStatsArray, proxyStats)

		}

	}

	/*

		SCRAPINGANT
	*/

	scrapingAntApiKeyList := []string{os.Getenv("SCRAPINGANT_INFO_API_KEY"), os.Getenv("SCRAPINGANT_IAN_API_KEY")}
	for _, scrapingAntApiKey := range scrapingAntApiKeyList {

		scrapeAccountEndpoint := "https://api.scrapingant.com/v1/usage?x-api-key=" + scrapingAntApiKey

		req, _ = http.NewRequest("GET", scrapeAccountEndpoint, nil)

		client = &http.Client{
			Timeout: 130 * time.Second,
		}

		// Make Request
		resp, err = client.Do(req)

		if err == nil && resp.StatusCode == 200 {

			defer resp.Body.Close()

			// Convert Body To JSON
			var scrapingAntAccount ScrapingAntAccount
			json.NewDecoder(resp.Body).Decode(&scrapingAntAccount)

			var subAccount string
			if scrapingAntApiKey == "6406a4c331874e9e90e0006a2fcb3f19" {
				subAccount = "(info)"
			} else if scrapingAntApiKey == "d4db1ea2587e4f93be58ad0379f6895c" {
				subAccount = "(ian)"
			}

			// endDate, err := time.ParseInLocation("2023-11-08T00:00:00", scrapingAntAccount.EndDate, timeZoneLocation)
			// if err != nil {
			// 	log.Println(err)
			// }

			endDate, err := time.Parse("2006-01-02T15:04:05.999999", scrapingAntAccount.EndDate)
			if err != nil {
				log.Println("ScrapingAnt Error parsing end date:", err)
			}

			// log.Println("scrapingAntAccount", "")
			// log.Println("scrapingAntAccount.EndDate", scrapingAntAccount.EndDate)
			// log.Println("endDate", endDate)

			proxyStats := ProxyProviderStats{
				Name:             "ScrapingAnt" + subAccount,
				CreditLimit:      scrapingAntAccount.PlanTotalCredits,
				UsedCredits:      scrapingAntAccount.PlanTotalCredits - scrapingAntAccount.RemainedCredits,
				UsedPercentage:   float64(scrapingAntAccount.PlanTotalCredits-scrapingAntAccount.RemainedCredits) / float64(scrapingAntAccount.PlanTotalCredits),
				RenewalDate:      endDate,
				DaysUntilRenewal: int(endDate.Sub(now).Hours() / 24),
			}

			proxyProviderStatsArray = append(proxyProviderStatsArray, proxyStats)

		}

	}

	/*

		SCRAPE.DO

	*/

	scrapedoApiKey := os.Getenv("SCRAPEDO_API_KEY")
	scrapedoAccountEndpoint := "https://api.Scrape.do/info?token=" + scrapedoApiKey

	req, _ = http.NewRequest("GET", scrapedoAccountEndpoint, nil)

	client = &http.Client{
		Timeout: 130 * time.Second,
	}

	// Make Request
	resp, err = client.Do(req)

	if err == nil && resp.StatusCode == 200 {

		defer resp.Body.Close()

		// Convert Body To JSON
		var scrapedoAccount ScrapedoAccount
		json.NewDecoder(resp.Body).Decode(&scrapedoAccount)

		usedApiCredits := scrapedoAccount.MaxMonthlyRequest - scrapedoAccount.RemainingMonthlyRequest

		proxyStats := ProxyProviderStats{
			Name:             "ScrapeDo",
			CreditLimit:      scrapedoAccount.MaxMonthlyRequest,
			UsedCredits:      usedApiCredits,
			UsedPercentage:   float64(usedApiCredits) / float64(scrapedoAccount.MaxMonthlyRequest),
			DaysUntilRenewal: 999,
		}

		proxyProviderStatsArray = append(proxyProviderStatsArray, proxyStats)
	}

	/*

		INFATICA

	*/

	/*

		ZENROWS

	*/

	zenrowsApiKey := os.Getenv("ZENROWS_API_KEY")
	zenrowsAccountEndpoint := "https://api.zenrows.com/v1/usage?apikey=" + zenrowsApiKey

	req, _ = http.NewRequest("GET", zenrowsAccountEndpoint, nil)

	client = &http.Client{
		Timeout: 130 * time.Second,
	}

	// Make Request
	resp, err = client.Do(req)

	if err == nil && resp.StatusCode == 200 {

		defer resp.Body.Close()

		// Convert Body To JSON
		var zenrowsAccount ZenrowsAccount
		json.NewDecoder(resp.Body).Decode(&zenrowsAccount)

		proxyStats := ProxyProviderStats{
			Name:             "ZenRows",
			CreditLimit:      zenrowsAccount.ApiCreditLimit,
			UsedCredits:      zenrowsAccount.ApiCreditUsage,
			UsedPercentage:   float64(zenrowsAccount.ApiCreditUsage) / float64(zenrowsAccount.ApiCreditLimit),
			DaysUntilRenewal: 999,
		}

		proxyProviderStatsArray = append(proxyProviderStatsArray, proxyStats)
	}

	/*

		ZYTE

	*/

	// Get Start & End Date
	startDate := now
	endDate := now
	currentDay := now.Day()

	if currentDay < 16 {
		startDate = time.Date(now.Year(), now.Month()-1, 17, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(now.Year(), now.Month(), 17, 0, 0, 0, 0, time.UTC)
	} else {
		startDate = time.Date(now.Year(), now.Month(), 17, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(now.Year(), now.Month()+1, 17, 0, 0, 0, 0, time.UTC)
	}

	startDateString := startDate.Format("2006-01-02T15:04")
	endDateString := endDate.Format("2006-01-02T15:04")

	zyteApiKey := os.Getenv("ZYTE_SPM_API_KEY")
	zyteAccountEndpoint := "https://crawlera-stats.scrapinghub.com/stats/?start_date=" + startDateString + "&end_date=" + endDateString
	zyteRequestLimit := 30000000

	req, _ = http.NewRequest("GET", zyteAccountEndpoint, nil)

	client = &http.Client{
		Timeout: 130 * time.Second,
	}

	// Add Headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(zyteApiKey, "")

	// Make Request
	resp, err = client.Do(req)

	if err == nil && resp.StatusCode == 200 {

		defer resp.Body.Close()

		// Convert Body To JSON
		var zyteResponse ZyteResponse
		json.NewDecoder(resp.Body).Decode(&zyteResponse)

		// log.Println("zyteResponse", zyteResponse)

		if len(zyteResponse.Results) > 0 {

			usedApiCredits := zyteResponse.Results[0].Clean

			proxyStats := ProxyProviderStats{
				Name:             "Zyte_SPM",
				CreditLimit:      zyteRequestLimit,
				UsedCredits:      usedApiCredits,
				UsedPercentage:   float64(usedApiCredits) / float64(zyteRequestLimit),
				RenewalDate:      endDate,
				DaysUntilRenewal: int(endDate.Sub(now).Hours() / 24),
			}

			proxyProviderStatsArray = append(proxyProviderStatsArray, proxyStats)

		}

	}

	/*

		SORT STATS BY USAGE

	*/

	sort.SliceStable(proxyProviderStatsArray, func(i, j int) bool {
		return proxyProviderStatsArray[i].UsedPercentage > proxyProviderStatsArray[j].UsedPercentage
	})

	/*

		SEND SLACK MESSAGE

	*/

	// Create Stats Blob
	statsString := "```Proxy | Limit | Used | % | Days |\n"
	for _, proxyProviderStat := range proxyProviderStatsArray {
		statsString = statsString + proxyProviderStat.Name + " | "

		var creditLimit string
		if proxyProviderStat.CreditLimit >= 1000000 {
			millions := float64(proxyProviderStat.CreditLimit) / 1000000
			creditLimit = fmt.Sprintf("%v", math.Round(millions*100)/100) + "M"
		} else if proxyProviderStat.CreditLimit < 1000000 {
			thousands := float64(proxyProviderStat.CreditLimit) / 1000
			creditLimit = fmt.Sprintf("%v", math.Round(thousands*100)/100) + "k"
		} else {
			creditLimit = fmt.Sprintf("%v", proxyProviderStat.CreditLimit)
		}

		var usedCredits string
		if proxyProviderStat.UsedCredits >= 1000000 {
			millions := float64(proxyProviderStat.UsedCredits) / 1000000
			usedCredits = fmt.Sprintf("%v", math.Round(millions*100)/100) + "M"
		} else if proxyProviderStat.UsedCredits < 1000 {
			usedCredits = fmt.Sprintf("%v", proxyProviderStat.UsedCredits)
		} else if proxyProviderStat.UsedCredits < 1000000 {
			thousands := float64(proxyProviderStat.UsedCredits) / 1000
			usedCredits = fmt.Sprintf("%v", math.Round(thousands*100)/100) + "k"
		} else {
			usedCredits = fmt.Sprintf("%v", proxyProviderStat.UsedCredits)
		}

		statsString = statsString + creditLimit + " | "
		statsString = statsString + usedCredits + " | "
		statsString = statsString + fmt.Sprintf("%v", (math.Round(proxyProviderStat.UsedPercentage*100)/100)*100) + "% | "
		if proxyProviderStat.DaysUntilRenewal == 999 {
			statsString = statsString + "- |\n"
		} else {
			statsString = statsString + fmt.Sprintf("%v", proxyProviderStat.DaysUntilRenewal) + " |\n"
		}

	}
	statsString = statsString + "Infatica | No API Endpoint |\n"
	// statsString = statsString + "Zyte SPM | No API Endpoint |\n"
	statsString = statsString + "Zyte SB | No API Endpoint |\n"
	statsString = statsString + "Scrapingfish | Not implemented yet |\n"
	statsString = statsString + "```"

	// Send Slack Message
	headline := "Proxy Account Stats: " + dayStartDateString
	slack.SlackStatsAlert("#proxy-provider-accounts", headline, statsString)

	// for _, proxyProviderStat := range proxyProviderStatsArray {

	// 	log.Println("", "")
	// 	log.Println("Name:", proxyProviderStat.Name)
	// 	log.Println("CreditLimit:", proxyProviderStat.CreditLimit)
	// 	log.Println("UsedCredits:", proxyProviderStat.UsedCredits)
	// 	log.Println("UsedPercentage:", proxyProviderStat.UsedPercentage)
	// 	log.Println("RenewalDate:", proxyProviderStat.RenewalDate)
	// 	log.Println("DaysUntilRenewal:", proxyProviderStat.DaysUntilRenewal)
	// 	log.Println("", "")

	// }

	// log.Println("proxyProviderStatsArray", proxyProviderStatsArray)

}
