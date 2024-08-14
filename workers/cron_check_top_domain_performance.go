package workers

import (
	"github.com/go-co-op/gocron"
	"go_proxy_worker/logger"
	"go_proxy_worker/slack"
	// "go_proxy_worker/utils"
	"go_proxy_worker/db"
	// "strconv"
	// "strings"
	"sort"
	"math"
	"time"
	"fmt"
	// "log"

)




type DomainPerformance struct {
	// Day							time.Time           `json:"day"`
	Domain						string           `json:"domain"`
	Requests					uint           `json:"requests"`
	Success						uint           `json:"success"`
	Failed						uint           `json:"failed"`
	SuccessRate					float64        `json:"success_rate"`
}


func CronCheckTopDomainPerformance() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(12).Hours().Do(CheckTopDomainPerformance)
	s.StartBlocking()
}



func CheckTopDomainPerformance(){

	fileName := "cron_check_top_domain_performance.go"

	emptyErrMap := make(map[string]interface{})

	// load DB
	var db = db.GetDB()

	// Alert Thresholds
	// var successRateThreshold float64 = 90

	// Get Today's Date
	now := time.Now()
	yesterday := now.Add(-time.Hour * time.Duration(24))
	dayStartTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	dayStartDateString := dayStartTime.Format("2006-01-02 15:04")
	dayStartDate, _ := time.Parse("2006-01-02 15:04", dayStartDateString)
	// log.Println("dayStartDate.UTC()", dayStartDate.UTC())

	// Top Domain Stats
	var domainPerformanceArray []DomainPerformance
	domainPerformanceResult := db.Raw(`
	SELECT 
	account_proxy_domain_stat_domain as domain,
	account_proxy_domain_stat_requests as requests,
	account_proxy_domain_stat_successful as success,
	account_proxy_domain_stat_failed as failed,
	(account_proxy_domain_stat_successful * 1.0/account_proxy_domain_stat_requests)*100 as success_rate
	from account_proxy_domain_stats apds 
	where account_proxy_domain_stat_day_start_time > ? and account_proxy_domain_stat_requests > 0
	order by requests desc
	limit 30
	`, dayStartDate.UTC()).Scan(&domainPerformanceArray)

	if domainPerformanceResult.Error != nil {
		logger.LogError("INFO", fileName, domainPerformanceResult.Error, "failed to get DomainPerformance from DB", emptyErrMap)
	}


	/*
	
		SORT STATS BY USAGE

	*/

	sort.SliceStable(domainPerformanceArray, func(i, j int) bool {
		return domainPerformanceArray[i].SuccessRate < domainPerformanceArray[j].SuccessRate
	})
	
	/*
	
		SEND SLACK MESSAGE

	*/


	if len(domainPerformanceArray) > 0 {

		// Send Slack Message
		statsString := "```Domain | Req | Success | % |\n" 
		for _, domainPerformanceStat := range domainPerformanceArray {

			var requests string
			if domainPerformanceStat.Requests >= 1000000 {
				millions := float64(domainPerformanceStat.Requests) / 1000000
				requests = fmt.Sprintf("%v", math.Round(millions*100)/100) + "M"
			} else if domainPerformanceStat.Requests < 1000000 && domainPerformanceStat.Requests >= 1000 {
				thousands := float64(domainPerformanceStat.Requests) / 1000
				requests = fmt.Sprintf("%v", math.Round(thousands*100)/100) + "k"
			} else {
				requests = fmt.Sprintf("%v", domainPerformanceStat.Requests)
			}

			var successful string
			if domainPerformanceStat.Success >= 1000000 {
				millions := float64(domainPerformanceStat.Success) / 1000000
				successful = fmt.Sprintf("%v", math.Round(millions*100)/100) + "M"
			} else if domainPerformanceStat.Success < 1000000 && domainPerformanceStat.Success >= 1000 {
				thousands := float64(domainPerformanceStat.Success) / 1000
				successful = fmt.Sprintf("%v", math.Round(thousands*100)/100) + "k"
			} else {
				successful = fmt.Sprintf("%v", domainPerformanceStat.Success)
			}

			// var failed string
			// if domainPerformanceStat.Failed >= 1000000 {
			// 	millions := float64(domainPerformanceStat.Failed) / 1000000
			// 	failed = fmt.Sprintf("%v", math.Round(millions*100)/100) + "M"
			// } else if domainPerformanceStat.Failed < 1000000 && domainPerformanceStat.Failed >= 1000 {
			// 	thousands := float64(domainPerformanceStat.Failed) / 1000
			// 	failed = fmt.Sprintf("%v", math.Round(thousands*100)/100) + "k"
			// } else {
			// 	failed = fmt.Sprintf("%v", domainPerformanceStat.Failed)
			// }


			statsString = statsString + domainPerformanceStat.Domain + " | "
			statsString = statsString + requests + " | "
			statsString = statsString + successful + " | "
			// statsString = statsString + failed + " | "
			statsString = statsString + fmt.Sprintf("%v", int(domainPerformanceStat.SuccessRate)) + "% |\n"
		}
		statsString = statsString + "```"

		headline := "Domain Performance Stats: " + dayStartDateString
		slack.SlackStatsAlert("#domain-performance-stats", headline, statsString)

	}



}