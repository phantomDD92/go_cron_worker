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


type EnterpriseUser struct {
	AccountId					uint           `json:"account_id"`
	AccountName					string         `json:"account_name"`
	ProxyPlanName				string         `json:"proxy_plan_name"`
	ProxyPlanPrice				uint           `json:"proxy_plan_price"`
	ApiCreditLimit				uint           `json:"api_credit_limit"`
}


type UserPerformance struct {
	// Day							time.Time           `json:"day"`
	Domain						string           `json:"domain"`
	Requests					uint           `json:"requests"`
	Success						uint           `json:"success"`
	Failed						uint           `json:"failed"`
	SuccessRate					float64        `json:"success_rate"`
}


func CronCheckEnterpriseUserPerformance() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(12).Hours().Do(CheckEnterpriseUserPerformance)
	s.StartBlocking()
}



func CheckEnterpriseUserPerformance(){

	fileName := "cron_check_enterprise_user_performance.go"

	emptyErrMap := make(map[string]interface{})

	// load DB
	var db = db.GetDB()

	// Alert Thresholds
	var successRateThreshold float64 = 90

	// Get Today's Date
	now := time.Now()
	yesterday := now.Add(-time.Hour * time.Duration(24))
	dayStartTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	dayStartDateString := dayStartTime.Format("2006-01-02 15:04")
	dayStartDate, _ := time.Parse("2006-01-02 15:04", dayStartDateString)
	// log.Println("dayStartDate.UTC()", dayStartDate.UTC())

	// List of Enterprise Users
	var enterpriseUserArray []EnterpriseUser
	enterpriseUserResult := db.Raw(`
	select 
	accounts.id as account_id,
	accounts.name as account_name,
	proxy_plans.name as proxy_plan_name,
	proxy_plans.price as proxy_plan_price,
	proxy_plans.api_credit_limit as api_credit_limit
	from accounts accounts
	join proxy_plans proxy_plans 
	on accounts.proxy_plan_id = proxy_plans.id 
	where proxy_plans.price > 250
	`).Scan(&enterpriseUserArray)

	if enterpriseUserResult.Error != nil {
		logger.LogError("INFO", fileName, enterpriseUserResult.Error, "failed to get EnterpriseUsers from DB", emptyErrMap)
	}


	for _, enterpriseUser := range enterpriseUserArray {

		// log.Println("proxyProvider", proxyProvider)

		// Get Performance Stats
		var userPerformanceArray []UserPerformance
		userPerformanceResult := db.Raw(`
		SELECT 
		account_proxy_domain_stat_domain as domain,
		account_proxy_domain_stat_requests as requests,
		account_proxy_domain_stat_successful as success,
		account_proxy_domain_stat_failed as failed,
		(account_proxy_domain_stat_successful * 1.0/account_proxy_domain_stat_requests)*100 as success_rate
		from account_proxy_domain_stats  
		where account_id = ? and account_proxy_domain_stat_day_start_time > ?
		`, enterpriseUser.AccountId, dayStartDate.UTC()).Scan(&userPerformanceArray)

		if userPerformanceResult.Error != nil {
			logger.LogError("INFO", fileName, userPerformanceResult.Error, "failed to get UserPerformance from DB", emptyErrMap)
		}

		
		var lowPerformanceArray []UserPerformance
		for _, userPerformance := range userPerformanceArray {
			if userPerformance.SuccessRate < successRateThreshold {
				lowPerformanceArray = append(lowPerformanceArray, userPerformance)
			}
		}

		// log.Println("proxyProviderStatsArray", proxyProviderStatsArray)
		// log.Println("failedValidationProxyStatsArray", failedValidationProxyStatsArray)

		if len(lowPerformanceArray) > 0 {

			// Send Slack Message
			statsString := "```Domain | Req | Success | Fail | % |\n" 
			for _, lowPerformanceStat := range lowPerformanceArray {

				var requests string
				if lowPerformanceStat.Requests >= 1000000 {
					millions := float64(lowPerformanceStat.Requests) / 1000000
					requests = fmt.Sprintf("%v", math.Round(millions*100)/100) + "M"
				} else if lowPerformanceStat.Requests < 1000000 && lowPerformanceStat.Requests >= 1000 {
					thousands := float64(lowPerformanceStat.Requests) / 1000
					requests = fmt.Sprintf("%v", math.Round(thousands*100)/100) + "k"
				} else {
					requests = fmt.Sprintf("%v", lowPerformanceStat.Requests)
				}

				var successful string
				if lowPerformanceStat.Success >= 1000000 {
					millions := float64(lowPerformanceStat.Success) / 1000000
					successful = fmt.Sprintf("%v", math.Round(millions*100)/100) + "M"
				} else if lowPerformanceStat.Success < 1000000 && lowPerformanceStat.Success >= 1000 {
					thousands := float64(lowPerformanceStat.Success) / 1000
					successful = fmt.Sprintf("%v", math.Round(thousands*100)/100) + "k"
				} else {
					successful = fmt.Sprintf("%v", lowPerformanceStat.Success)
				}

				var failed string
				if lowPerformanceStat.Failed >= 1000000 {
					millions := float64(lowPerformanceStat.Failed) / 1000000
					failed = fmt.Sprintf("%v", math.Round(millions*100)/100) + "M"
				} else if lowPerformanceStat.Failed < 1000000 && lowPerformanceStat.Failed >= 1000 {
					thousands := float64(lowPerformanceStat.Failed) / 1000
					failed = fmt.Sprintf("%v", math.Round(thousands*100)/100) + "k"
				} else {
					failed = fmt.Sprintf("%v", lowPerformanceStat.Failed)
				}


				statsString = statsString + lowPerformanceStat.Domain + " | "
				statsString = statsString + requests + " | "
				statsString = statsString + successful + " | "
				statsString = statsString + failed + " | "
				statsString = statsString + fmt.Sprintf("%v", int(lowPerformanceStat.SuccessRate)) + "% |\n"
			}
			statsString = statsString + "```"

			headline := strings.Title(enterpriseUser.AccountName) + " (" + fmt.Sprintf("%v", enterpriseUser.AccountId) + "): $" + fmt.Sprintf("%v", enterpriseUser.ProxyPlanPrice)
			slack.SlackStatsAlert("#enterprise-performance-stats", headline, statsString)

		}

	}





}