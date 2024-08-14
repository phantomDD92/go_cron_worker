package workers

import (
	"github.com/go-co-op/gocron"
	"go_proxy_worker/logger"
	"go_proxy_worker/slack"
	// "go_proxy_worker/utils"
	"go_proxy_worker/db"
	// "strconv"
	"strings"
	"math"
	"time"
	"fmt"
	// "log"

)


type ProxyProviderStat struct {
	Day							time.Time           `json:"day"`
	Domain						string           `json:"domain"`
	Requests					uint           `json:"requests"`
	ProxySuccessful				uint           `json:"proxy_successful"`
	FailedValidation			uint           `json:"failed_validation"`
	InvalidRate					float64        `json:"invalid_rate"`
}




func CronCheckProxyProviderFailedValidation() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(12).Hours().Do(CheckProxyProviderFailedValidation)
	s.StartBlocking()
}



func CheckProxyProviderFailedValidation(){

	fileName := "cron_check_proxy_provider_failed_validation.go"

	emptyErrMap := make(map[string]interface{})

	// load DB
	var db = db.GetDB()

	// Alert Thresholds
	var numberFailedValidationPages uint = 1000
	var percentageFailedValidationPages float64 = 10

	// Get Today's Date
	now := time.Now()
	yesterday := now.Add(-time.Hour * time.Duration(24))
	dayStartTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	dayStartDateString := dayStartTime.Format("2006-01-02 15:04")
	dayStartDate, _ := time.Parse("2006-01-02 15:04", dayStartDateString)
	// log.Println("dayStartDate.UTC()", dayStartDate.UTC())

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

		// log.Println("proxyProvider", proxyProvider)

		// Get Proxy Provider Stats
		var proxyProviderStatsArray []ProxyProviderStat
		proxyProviderStatsResult := db.Raw(`
		select 
		sops_day_proxy_domain_stat_day_start_time as day,
		sops_day_proxy_domain_stat_domain as domain,
		sops_day_proxy_domain_stat_requests as requests,
		sops_day_proxy_domain_stat_successful + sops_day_proxy_domain_stat_failed_validation as proxy_successful,
		sops_day_proxy_domain_stat_failed_validation as failed_validation, 
		(sops_day_proxy_domain_stat_failed_validation*1.0 / (sops_day_proxy_domain_stat_successful + sops_day_proxy_domain_stat_failed_validation)*1.0)*100 as invalid_rate
		from sops_day_proxy_domain_stats
		where sops_proxy_name = ? and sops_day_proxy_domain_stat_day_start_time > ? and sops_day_proxy_domain_stat_successful >= 1
		order by invalid_rate desc
		`, proxyProvider, dayStartDate.UTC()).Scan(&proxyProviderStatsArray)

		if proxyProviderStatsResult.Error != nil {
			logger.LogError("INFO", fileName, proxyProviderStatsResult.Error, "failed to get ProxyProviderStats from DB", emptyErrMap)
		}

		

		var failedValidationProxyStatsArray []ProxyProviderStat
		for _, proxyProviderStat := range proxyProviderStatsArray {
			if proxyProviderStat.FailedValidation > numberFailedValidationPages || (proxyProviderStat.InvalidRate > percentageFailedValidationPages && proxyProviderStat.FailedValidation > 100) {
				failedValidationProxyStatsArray = append(failedValidationProxyStatsArray, proxyProviderStat)
			}
		}

		// log.Println("proxyProviderStatsArray", proxyProviderStatsArray)
		// log.Println("failedValidationProxyStatsArray", failedValidationProxyStatsArray)

		if len(failedValidationProxyStatsArray) > 0 {

			// Send Slack Message
			statsString := "```Domain | Req | Success | Invalid | % |\n" 
			for _, failedValidationProxyStat := range failedValidationProxyStatsArray {
				statsString = statsString + failedValidationProxyStat.Domain + " | "
				statsString = statsString + fmt.Sprintf("%v", failedValidationProxyStat.Requests) + " | "
				statsString = statsString + fmt.Sprintf("%v", failedValidationProxyStat.ProxySuccessful) + " | "
				statsString = statsString + fmt.Sprintf("%v", failedValidationProxyStat.FailedValidation) + " | "
				statsString = statsString + fmt.Sprintf("%v", math.Round(failedValidationProxyStat.InvalidRate*100)/100) + "% |\n"
			}
			statsString = statsString + "```"

			headline := strings.Title(proxyProvider) + ": Failed Validation Stats"
			slack.SlackStatsAlert("#proxy-provider-failed-validation", headline, statsString)

		}
		

	}





}