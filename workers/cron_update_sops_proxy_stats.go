package workers

import (
	"github.com/go-co-op/gocron"
	"go_proxy_worker/models"
	"go_proxy_worker/logger"
	"go_proxy_worker/utils"
	"go_proxy_worker/dbRedisQueries"
	"github.com/fatih/structs"
	"go_proxy_worker/db"
	// "gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	"log"
	"fmt"
)





func CronUpdateProxyStats() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Minutes().Do(RunUpdateProxyStats)
	s.StartBlocking()
}


func RunUpdateProxyStats(){

	// go NEWRunUpdateProxyStats()

	fileName := "cron_update_sops_proxy_stats.go"

	//emptyErrMap := make(map[string]interface{})

	// Redis Details
	// var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	// load DB
	var db = db.GetDB()


	// Get Active Proxy Accounts
	activeProxyKeySet := "activeProxysKeySet"
	// activeProxyKeySet = utils.RedisEnvironmentVersion(activeProxyKeySet)

	listActiveProxyKeys, err := statsProxyRedisClient.SMembers(redisContext, activeProxyKeySet).Result()
	if err != nil {
		errData := map[string]interface{}{
			"activeProxyKeySet": activeProxyKeySet,
		}
		logger.LogError("INFO", fileName, err, "activeProxyKeySet not in Redis", errData)
	}


	log.Println("listActiveProxyKeys", listActiveProxyKeys)

	for _, proxy := range listActiveProxyKeys {

		log.Println("")
		log.Println("###################")
		log.Println("")
		log.Println("proxy", proxy)
		log.Println("")
		log.Println("")


		// Remove From activeAccountAccountsKeySet
		err = statsProxyRedisClient.SRem(redisContext, activeProxyKeySet, proxy).Err()
		if err != nil {
			errData := map[string]interface{}{
				"proxy": proxy,
				"activeProxyKeySet": activeProxyKeySet,
			}
			logger.LogError("ERROR", fileName, err, "failed to delete proxy from activeProxyKeySet in Redis", errData)
		}


		dayProxyStatsMap := make(map[uint]models.SopsDayProxyStat)
		dayProxyDomainStatsMap := make(map[uint]map[string]models.SopsDayProxyDomainStat)

		// Get Active Time Windows For Proxy
		proxyTimeWindowKeySet := "proxyTimeWindowKeySet?proxyName=" + proxy
		//proxyTimeWindowKeySet = utils.RedisEnvironmentVersion(proxyTimeWindowKeySet)
		listActiveTimeWindowsKeys, err := statsProxyRedisClient.SMembers(redisContext, proxyTimeWindowKeySet).Result()
		if err != nil {
			errData := map[string]interface{}{
				"proxyTimeWindowKeySet": proxyTimeWindowKeySet,
			}
			logger.LogError("INFO", fileName, err, "proxyTimeWindowKeySet not in Redis", errData)
		}


		log.Println("listActiveTimeWindowsKeys", listActiveTimeWindowsKeys)


		sopsProxyProviderId := dbRedisQueries.GetSopsProxyProviderId(proxy, db, statsProxyRedisClient, redisContext, fileName)



		/*
			FOR EACH TIME WINDOW GET THE FOLLOWING FROM REDIS

			- Total requests
			- Successful requests
			- Failed requests
			- API Credits

			And add to the totals for that day

		*/


		var key string

		for _, timeWindowRaw := range listActiveTimeWindowsKeys {

			log.Println("")
			log.Println("timeWindowRaw", timeWindowRaw)
			log.Println("")


			
			timeWindow := strings.Trim(timeWindowRaw, `"`)
			splitTimeWindow := strings.Split(timeWindow, "::")

			// Date
			splitDate := strings.Split(splitTimeWindow[0], "-") 
			year, _ := strconv.Atoi(splitDate[0])
			day, _ := strconv.Atoi(splitDate[2])
			month := utils.ConvertMonthStringToInt(splitDate[1])

			// Time
			// splitTime := strings.Split(splitTimeWindow[1], "-") 
			// hour, _ := strconv.Atoi(splitTime[0])

			dayStartTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			nowUTC := time.Now().UTC()


			sopsDayProxyStatId := dbRedisQueries.GetSopsDayProxyStatId(sopsProxyProviderId, proxy, dayStartTime, dayStartTime.String(), db, statsProxyRedisClient, redisContext, fileName)

			if _, ok := dayProxyStatsMap[sopsDayProxyStatId]; !ok {
				dayProxyStatsMap[sopsDayProxyStatId] = models.SopsDayProxyStat{}
			}

			dayProxyStats := dayProxyStatsMap[sopsDayProxyStatId]
			if dayProxyStats.ID == 0 {
				dayProxyStats.ID = sopsDayProxyStatId
			}

			dayProxyStats.SopsProxyName = proxy
			dayProxyStats.SopsProxyProviderId = sopsProxyProviderId
			dayProxyStats.SopsDayProxyStatDayStartTime = dayStartTime

			// Get Total Requests
			key = "proxyTotalRequests?proxyName=" + proxy + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
			// key = utils.RedisEnvironmentVersion(key)
			dayProxyStats.SopsDayProxyStatRequests = dayProxyStats.SopsDayProxyStatRequests + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


			// Get Total Successful Requests
			key = "proxySuccessfulRequests?proxyName=" + proxy + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
			// key = utils.RedisEnvironmentVersion(key)
			dayProxyStats.SopsDayProxyStatSuccessful = dayProxyStats.SopsDayProxyStatSuccessful + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


			// Get Total Successful + Failed Validation Requests
			key = "proxySuccessfulFailedValidationRequests?proxyName=" + proxy + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
			// key = utils.RedisEnvironmentVersion(key)
			dayProxyStats.SopsDayProxyStatFailedValidation = dayProxyStats.SopsDayProxyStatFailedValidation + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)

			// Get Total Failed Requests
			if dayProxyStats.SopsDayProxyStatRequests >= dayProxyStats.SopsDayProxyStatSuccessful {
				dayProxyStats.SopsDayProxyStatFailed = (dayProxyStats.SopsDayProxyStatRequests - dayProxyStats.SopsDayProxyStatSuccessful)
			}
			


			// Get Total API Credits
			key = "proxyUsedApiCredits?proxyName=" + proxy + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
			// key = utils.RedisEnvironmentVersion(key)
			dayProxyStats.SopsDayProxyStatCredits = dayProxyStats.SopsDayProxyStatCredits + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


			// Get Total Failed API Credits
			key = "proxyUsedFailedValidationApiCredits?proxyName=" + proxy + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
			// key = utils.RedisEnvironmentVersion(key)
			dayProxyStats.SopsDayProxyStatCreditsFailedValidation = dayProxyStats.SopsDayProxyStatCreditsFailedValidation + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


			// Get Total Latency Requests
			key = "proxyTotalLatency?proxyName=" + proxy + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
			// key = utils.RedisEnvironmentVersion(key)
			dayProxyStats.SopsDayProxyStatLatency = dayProxyStats.SopsDayProxyStatLatency + (dbRedisQueries.GetFloat64Redis(key, statsProxyRedisClient, redisContext, fileName))/1000


			dayProxyStatsMap[sopsDayProxyStatId] = dayProxyStats


			/*
				FOR EACH DOMAIN GET THE FOLLOWING FROM REDIS

				- Total requests
				- Successful requests
				- Failed requests
				- API Credits

				And add to the totals for that day

			*/

			proxyDomainKeySet := "proxyDomainKeySet?proxyName=" + proxy + "&timeWindow=" + fmt.Sprintf("%v", timeWindow)
			//proxyDomainKeySet = utils.RedisEnvironmentVersion(proxyDomainKeySet)
			listActiveDomainsKeys, err := statsProxyRedisClient.SMembers(redisContext, proxyDomainKeySet).Result()
			if err != nil {
				errData := map[string]interface{}{
					"listActiveDomainsKeys": listActiveDomainsKeys,
				}
				logger.LogError("INFO", fileName, err, "listActiveDomainsKeys not in Redis", errData)
			}

			log.Println("listActiveDomainsKeys", listActiveDomainsKeys)


			for _, domain := range listActiveDomainsKeys {

				log.Println("")
				log.Println("domain", timeWindowRaw)
				log.Println("")

				domain = strings.Trim(domain, `"`)

				if _, ok := dayProxyDomainStatsMap[sopsDayProxyStatId]; !ok {
					temp := make(map[string]models.SopsDayProxyDomainStat)
					temp[domain] = models.SopsDayProxyDomainStat{}
					dayProxyDomainStatsMap[sopsDayProxyStatId] = temp
				}


				domainProxyStats := dayProxyDomainStatsMap[sopsDayProxyStatId][domain]
				domainProxyStats.SopsProxyName = proxy
				domainProxyStats.SopsDayProxyStatId = sopsDayProxyStatId
				domainProxyStats.SopsProxyProviderId = sopsProxyProviderId
				domainProxyStats.SopsDayProxyDomainStatDomain = domain
				domainProxyStats.SopsDayProxyDomainStatDayStartTime = dayStartTime


				// Get Total Requests
				key = "proxyDomainRequests?proxyName=" + proxy + "&domain=" + domain + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
				// key = utils.RedisEnvironmentVersion(key)
				domainProxyStats.SopsDayProxyDomainStatRequests = domainProxyStats.SopsDayProxyDomainStatRequests + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)

				// Get Total Successful Requests
				key = "proxySuccessfulRequests?proxyName=" + proxy + "&domain=" + domain + "&timeWindow=" + fmt.Sprintf("%v", timeWindow)
				// key = utils.RedisEnvironmentVersion(key)
				domainProxyStats.SopsDayProxyDomainStatSuccessful = domainProxyStats.SopsDayProxyDomainStatSuccessful + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


				// Get Total Successful + Failed Validation Requests
				key = "proxySuccessfulFailedValidationRequests?proxyName=" + proxy + "&domain=" + domain + "&timeWindow=" + fmt.Sprintf("%v", timeWindow)
				// key = utils.RedisEnvironmentVersion(key)
				domainProxyStats.SopsDayProxyDomainStatFailedValidation = domainProxyStats.SopsDayProxyDomainStatFailedValidation + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)



				// Get Total Failed Requests
				if domainProxyStats.SopsDayProxyDomainStatRequests >= domainProxyStats.SopsDayProxyDomainStatSuccessful {
					domainProxyStats.SopsDayProxyDomainStatFailed = (domainProxyStats.SopsDayProxyDomainStatRequests - domainProxyStats.SopsDayProxyDomainStatSuccessful)
				}
				


				// Get Total API Credits
				key = "proxyUsedApiCredits?proxyName=" + proxy + "&domain=" + domain + "&timeWindow=" + fmt.Sprintf("%v", timeWindow)
				// key = utils.RedisEnvironmentVersion(key)
				domainProxyStats.SopsDayProxyDomainStatCredits = domainProxyStats.SopsDayProxyDomainStatCredits + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


				// Get Total Failed API Credits
				key = "proxyUsedFailedValidationApiCredits?proxyName=" + proxy + "&domain=" + domain + "&timeWindow=" + fmt.Sprintf("%v", timeWindow)
				// key = utils.RedisEnvironmentVersion(key)
				domainProxyStats.SopsDayProxyDomainStatCreditsFailedValidation = domainProxyStats.SopsDayProxyDomainStatCreditsFailedValidation + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


				// Get Total Latency Requests
				key = "proxyDomainLatency?proxyName=" + proxy + "&domain=" + domain + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
				// key = utils.RedisEnvironmentVersion(key)
				domainProxyStats.SopsDayProxyDomainStatLatency = domainProxyStats.SopsDayProxyDomainStatLatency + (dbRedisQueries.GetFloat64Redis(key, statsProxyRedisClient, redisContext, fileName))/1000

				dayProxyDomainStatsMap[sopsDayProxyStatId][domain] = domainProxyStats

				
			}
			

			// Remove Old Time Windows From Yesterday
			if nowUTC.Day() != day {

				log.Println("removing timeWindow", timeWindowRaw)
				log.Println("")

				// Remove timeWindow From proxyTimeWindowKeySet
				err = statsProxyRedisClient.SRem(redisContext, proxyTimeWindowKeySet, timeWindowRaw).Err()
				if err != nil {
					errData := map[string]interface{}{
						"timeWindowRaw": timeWindowRaw,
						"proxyTimeWindowKeySet": proxyTimeWindowKeySet,
					}
					logger.LogError("ERROR", fileName, err, "failed to delete timeWindowRaw from proxyTimeWindowKeySet in Redis", errData)
				}

			}

		}


		log.Println("dayProxyStatsMap", dayProxyStatsMap)
		log.Println("dayProxyDomainStatsMap", dayProxyDomainStatsMap)


		/*

			UPDATE DB

		*/


		// Update AccountProxyStats
		for sopsDayProxyStatId, dayProxyStat := range dayProxyStatsMap {


			sopsDayProxyStatsUpdateMap := map[string]interface{}{
				"sops_day_proxy_stat_requests": dayProxyStat.SopsDayProxyStatRequests,
				"sops_day_proxy_stat_successful": dayProxyStat.SopsDayProxyStatSuccessful,
				// "sops_day_proxy_stat_failed": dayProxyStat.SopsDayProxyStatFailed,
				"sops_day_proxy_stat_failed_validation": dayProxyStat.SopsDayProxyStatFailedValidation,
				"sops_day_proxy_stat_credits": dayProxyStat.SopsDayProxyStatCredits,
				"sops_day_proxy_stat_credits_failed_validation": dayProxyStat.SopsDayProxyStatCreditsFailedValidation,
				"sops_day_proxy_stat_latency": dayProxyStat.SopsDayProxyStatLatency,
			}

			if dayProxyStat.SopsDayProxyStatFailed > 0 {
				sopsDayProxyStatsUpdateMap["sops_day_proxy_stat_failed"] = dayProxyStat.SopsDayProxyStatFailed
			}

			log.Println("")
			log.Println("sopsDayProxyStatsUpdateMap", sopsDayProxyStatsUpdateMap)
			log.Println("")

			result := db.Model(&dayProxyStat).Where("id = ?", sopsDayProxyStatId).Updates(sopsDayProxyStatsUpdateMap)
			if result.Error != nil || result.RowsAffected == 0 {
				createResult := db.Create(&dayProxyStat)
				if createResult.Error != nil {
					errData := structs.Map(dayProxyStat)
					logger.LogError("ERROR", fileName, err, "Failed to update or create dayProxyStat in DB", errData)
				}
			}

		}


		// Update AccountProxyDomainStats
		for sopsDayProxyStatId, submap := range dayProxyDomainStatsMap {
			for domain, domainProxyStat := range submap {

				sopsDayProxyDomainStatsUpdateMap := map[string]interface{}{
					"sops_day_proxy_domain_stat_requests": domainProxyStat.SopsDayProxyDomainStatRequests,
					"sops_day_proxy_domain_stat_successful": domainProxyStat.SopsDayProxyDomainStatSuccessful,
					"sops_day_proxy_domain_stat_failed": domainProxyStat.SopsDayProxyDomainStatFailed,
					"sops_day_proxy_domain_stat_failed_validation": domainProxyStat.SopsDayProxyDomainStatFailedValidation,
					"sops_day_proxy_domain_stat_credits": domainProxyStat.SopsDayProxyDomainStatCredits,
					"sops_day_proxy_domain_stat_credits_failed_validation": domainProxyStat.SopsDayProxyDomainStatCreditsFailedValidation,
					"sops_day_proxy_domain_stat_latency": domainProxyStat.SopsDayProxyDomainStatLatency,
				}

				if domainProxyStat.SopsDayProxyDomainStatFailed > 0 {
					sopsDayProxyDomainStatsUpdateMap["sops_day_proxy_domain_stat_failed"] = domainProxyStat.SopsDayProxyDomainStatFailed
				}

				log.Println("")
				log.Println("sopsDayProxyDomainStatsUpdateMap", sopsDayProxyDomainStatsUpdateMap)
				log.Println("")

				result := db.Model(&domainProxyStat).Where("sops_day_proxy_stat_id = ? and sops_day_proxy_domain_stat_domain = ?", sopsDayProxyStatId, domain).Updates(sopsDayProxyDomainStatsUpdateMap)
				if result.Error != nil || result.RowsAffected == 0 {
					createResult := db.Create(&domainProxyStat)
					if createResult.Error != nil {
						errData := structs.Map(domainProxyStat)
						logger.LogError("ERROR", fileName, err, "Failed to update or create domainProxyStat in DB", errData)
					}
				}

			}
			
		}
		
		


	}

}