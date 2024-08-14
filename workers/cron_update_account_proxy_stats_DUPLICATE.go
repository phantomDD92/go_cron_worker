package workers

import (
	// "github.com/go-co-op/gocron"
	"go_proxy_worker/models"
	"go_proxy_worker/logger"
	"go_proxy_worker/utils"
	"go_proxy_worker/dbRedisQueries"
	"github.com/go-redis/redis/v8"
	// "github.com/fatih/structs"
	"go_proxy_worker/db"
	"gorm.io/gorm"
	"context"
	"strconv"
	"strings"
	"sync"
	"time"
	"log"
	"fmt"
)


func NEWRunAccountProxyStatsAggregation(dayString string, fileName string){

	//emptyErrMap := make(map[string]interface{})

	// Redis Details
	// var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	// load DB
	var db = db.GetDB()

	// Get Active Proxy Accounts
	activeAccountsKeySet := "activeAccountsKeySet?day=" + fmt.Sprintf("%v", dayString)
	// activeAccountsKeySet = utils.RedisEnvironmentVersion(activeAccountsKeySet)

	listActiveAccountsKeys, err := statsProxyRedisClient.SMembers(redisContext, activeAccountsKeySet).Result()
	if err != nil {
		errData := map[string]interface{}{
			"activeAccountsKeySet": activeAccountsKeySet,
		}
		logger.LogError("INFO", fileName, err, "activeAccountsKeySet not in Redis", errData)
	}


	log.Println("listActiveAccountsKeys", listActiveAccountsKeys)


	var wg sync.WaitGroup

	for _, accountIdString := range listActiveAccountsKeys {

		wg.Add(1)
		go NEWProcessAccountId(&wg, accountIdString, activeAccountsKeySet, dayString, statsProxyRedisClient, redisContext, db, fileName)

	}

	wg.Wait()

}


func NEWProcessAccountId(wg *sync.WaitGroup, accountIdString string, activeAccountsKeySet string, dayString string, statsProxyRedisClient *redis.Client, redisContext context.Context, db *gorm.DB, fileName string){


	defer wg.Done()


	activeProxyStatProcessingKey := "activeProcessingProxyStats?account_id=" + fmt.Sprintf("%v", accountIdString) 
	// activeProxyStatProcessingKey = utils.RedisEnvironmentVersion(activeProxyStatProcessingKey)

	// Check If Proxy Stats Is Active
	redisActiveProxyStatProcessing, _ := statsProxyRedisClient.Get(redisContext, activeProxyStatProcessingKey).Result()

	log.Println("activeProxyStatProcessingKey", activeProxyStatProcessingKey)
	log.Println("redisActiveProxyStatProcessing", redisActiveProxyStatProcessing)

	// If NOT Active --> Publish them to the Subscribers
	if redisActiveProxyStatProcessing == "" { 

		// Set Key As Being Processed
		err := statsProxyRedisClient.Set(redisContext, activeProxyStatProcessingKey, "true", 60*10*time.Second).Err() // 10 minute expiry
		if err != nil {
			errData := map[string]interface{}{
				"accountIdString": accountIdString,
				"activeProxyStatProcessingKey": activeProxyStatProcessingKey,
			}
			logger.LogError("WARN", fileName, err, "failed to set activeProxyStatProcessingKey in Redis", errData)
		}


		// Remove From activeAccountAccountsKeySet
		err = statsProxyRedisClient.SRem(redisContext, activeAccountsKeySet, accountIdString).Err()
		if err != nil {
			errData := map[string]interface{}{
				"accountIdString": accountIdString,
				"activeAccountsKeySet": activeAccountsKeySet,
			}
			logger.LogError("ERROR", fileName, err, "failed to delete accountIdString from activeAccountsKeySet in Redis", errData)
		}


		// Account Details
		accountId:= utils.StringToUint(accountIdString)
		accountItem, validAccount := dbRedisQueries.GetAccountDetails(accountId, db, statsProxyRedisClient, redisContext, fileName)

		if validAccount {

			accountProxyStatsMap := make(map[uint]models.AccountProxyStat)
			accountProxyDomainStatsMap := make(map[uint]map[string]models.AccountProxyDomainStat)

			// Get Active Time Windows For Account
			activeAccountTimewindowsKeySet := "activeAccountTimewindowsKeySet?account_id=" + fmt.Sprintf("%v", accountId) + "&day=" + fmt.Sprintf("%v", dayString)
			// activeAccountTimewindowsKeySet = utils.RedisEnvironmentVersion(activeAccountTimewindowsKeySet)
			listActiveTimeWindowsKeys, err := statsProxyRedisClient.SMembers(redisContext, activeAccountTimewindowsKeySet).Result()
			if err != nil {
				errData := map[string]interface{}{
					"activeAccountTimewindowsKeySet": activeAccountTimewindowsKeySet,
				}
				logger.LogError("INFO", fileName, err, "activeAccountTimewindowsKeySet not in Redis", errData)
			}



			log.Println("listActiveTimeWindowsKeys", listActiveTimeWindowsKeys)


			/*
				UPDATE API CREDITS USED
			*/

			var key string

			// // Get Stats From Redis
			// key = "accountRequests?account_id=" + fmt.Sprintf("%v", accountId)
			// accountProxyRequests :=  dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)

			// key = "accountSuccessfulRequests?account_id=" + fmt.Sprintf("%v", accountId)
			// accountProxySuccessful :=  dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)

			// key = "accountUsedApiCredits?account_id=" + fmt.Sprintf("%v", accountId)
			// accountProxyUsedCredits :=  dbRedisQueries.GetInt64Redis(key, coreProxyRedisClient, redisContext, fileName)


			/*
				FOR EACH TIME WINDOW GET THE FOLLOWING FROM REDIS

				- Total requests
				- Successful requests
				- Failed requests
				- API Credits

				And add to the totals for that day

			*/


			

			for _, timeWindowRaw := range listActiveTimeWindowsKeys {


				// // Remove From activeAccountAccountsKeySet
				// err := statsProxyRedisClient.SRem(redisContext, activeAccountTimewindowsKeySet, timeWindowRaw).Err()
				// if err != nil {
				// 	errData := map[string]interface{}{
				// 		"timeWindowRaw": timeWindowRaw,
				// 		"activeAccountTimewindowsKeySet": activeAccountTimewindowsKeySet,
				// 	}
				// 	logger.LogError("ERROR", fileName, err, "failed to delete timeWindowRaw from activeAccountTimewindowsKeySet in Redis", errData)
				// }

				
				timeWindow := strings.Trim(timeWindowRaw, `"`)
				splitTimeWindow := strings.Split(timeWindow, "::")

				// Date
				splitDate := strings.Split(splitTimeWindow[0], "-") 
				year, _ := strconv.Atoi(splitDate[0])
				day, _ := strconv.Atoi(splitDate[2])
				month := utils.ConvertMonthStringToInt(splitDate[1])

				// Time
				splitTime := strings.Split(splitTimeWindow[1], "-") 
				hour, _ := strconv.Atoi(splitTime[0])


				dateTime := time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC)
				usersDayStartTime := utils.GetUsersTimezoneDayStart(dateTime, accountItem) 
				usersDayStartTimeString := usersDayStartTime.String()

				accountProxyStatId :=  dbRedisQueries.GetAccountProxyStatId(accountId, usersDayStartTime, usersDayStartTimeString, db, statsProxyRedisClient, redisContext, fileName)

				if _, ok := accountProxyStatsMap[accountProxyStatId]; !ok {
					accountProxyStatsMap[accountProxyStatId] = models.AccountProxyStat{}
				}

				proxyStats := accountProxyStatsMap[accountProxyStatId]
				if proxyStats.ID == 0 {
					proxyStats.ID = accountProxyStatId
				}
				proxyStats.AccountId = accountId
				proxyStats.AccountProxyStatDayStartTime = usersDayStartTime

				// Get Total Requests
				key = "accountRequests?account_id=" + fmt.Sprintf("%v", accountId) + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
				// key = utils.RedisEnvironmentVersion(key)
				proxyStats.AccountProxyStatRequests = proxyStats.AccountProxyStatRequests + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


				// Get Total Successful Requests
				key = "accountSuccessfulRequests?account_id=" + fmt.Sprintf("%v", accountId) + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
				// key = utils.RedisEnvironmentVersion(key)
				proxyStats.AccountProxyStatSuccessful = proxyStats.AccountProxyStatSuccessful + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


				// Get Total Failed Requests
				proxyStats.AccountProxyStatFailed = (proxyStats.AccountProxyStatRequests - proxyStats.AccountProxyStatSuccessful)


				// Get Total API Credits
				key = "accountUsedApiCredits?account_id=" + fmt.Sprintf("%v", accountId) + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
				// key = utils.RedisEnvironmentVersion(key)
				proxyStats.AccountProxyStatCredits = proxyStats.AccountProxyStatCredits + dbRedisQueries.GetInt64Redis(key, statsProxyRedisClient, redisContext, fileName)


				accountProxyStatsMap[accountProxyStatId] = proxyStats


				/*
					FOR EACH DOMAIN GET THE FOLLOWING FROM REDIS

					- Total requests
					- Successful requests
					- Failed requests
					- API Credits

					And add to the totals for that day

				*/

				activeAccountDomainsKeySet := "activeAccountDomainsKeySet?account_id=" + fmt.Sprintf("%v", accountId) + "&day=" + fmt.Sprintf("%v", dayString)
				// activeAccountDomainsKeySet = utils.RedisEnvironmentVersion(activeAccountDomainsKeySet)
				listActiveDomainsKeys, err := statsProxyRedisClient.SMembers(redisContext, activeAccountDomainsKeySet).Result()
				if err != nil {
					errData := map[string]interface{}{
						"listActiveDomainsKeys": listActiveDomainsKeys,
					}
					logger.LogError("INFO", fileName, err, "listActiveDomainsKeys not in Redis", errData)
				}

				log.Println("listActiveDomainsKeys", listActiveDomainsKeys)


				for _, domain := range listActiveDomainsKeys {

					// // Remove From activeAccountAccountsKeySet
					// err := statsProxyRedisClient.SRem(redisContext, activeAccountDomainsKeySet, domain).Err()
					// if err != nil {
					// 	errData := map[string]interface{}{
					// 		"domain": domain,
					// 		"activeAccountDomainsKeySet": activeAccountDomainsKeySet,
					// 	}
					// 	logger.LogError("ERROR", fileName, err, "failed to delete domain from activeAccountDomainsKeySet in Redis", errData)
					// }



					domain = strings.Trim(domain, `"`)

					if _, ok := accountProxyDomainStatsMap[accountProxyStatId]; !ok {
						temp := make(map[string]models.AccountProxyDomainStat)
						temp[domain] = models.AccountProxyDomainStat{}
						accountProxyDomainStatsMap[accountProxyStatId] = temp
					}


					domainProxyStats := accountProxyDomainStatsMap[accountProxyStatId][domain]
					domainProxyStats.AccountId = accountId
					domainProxyStats.AccountProxyStatId = accountProxyStatId
					domainProxyStats.AccountProxyDomainStatDomain = domain
					domainProxyStats.AccountProxyDomainStatDayStartTime = usersDayStartTime


					// Get Total Requests
					key = "accountRequests?account_id=" + fmt.Sprintf("%v", accountId) + "&domain=" + domain + "&timeWindow=" + fmt.Sprintf("%v", timeWindow) 
					// key = utils.RedisEnvironmentVersion(key)
					domainProxyStats.AccountProxyDomainStatRequests = domainProxyStats.AccountProxyDomainStatRequests + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)

					// Get Total Successful Requests
					key = "accountSuccessfulRequests?account_id=" + fmt.Sprintf("%v", accountId) + "&domain=" + domain + "&timeWindow=" + fmt.Sprintf("%v", timeWindow)
					// key = utils.RedisEnvironmentVersion(key)
					domainProxyStats.AccountProxyDomainStatSuccessful = domainProxyStats.AccountProxyDomainStatSuccessful + dbRedisQueries.GetUintRedis(key, statsProxyRedisClient, redisContext, fileName)


					// Get Total Failed Requests
					domainProxyStats.AccountProxyDomainStatFailed = (domainProxyStats.AccountProxyDomainStatRequests - domainProxyStats.AccountProxyDomainStatSuccessful)


					// Get Total API Credits
					key = "accountUsedApiCredits?account_id=" + fmt.Sprintf("%v", accountId) + "&domain=" + domain + "&timeWindow=" + fmt.Sprintf("%v", timeWindow)
					// key = utils.RedisEnvironmentVersion(key)
					domainProxyStats.AccountProxyDomainStatCredits = domainProxyStats.AccountProxyDomainStatCredits + dbRedisQueries.GetInt64Redis(key, statsProxyRedisClient, redisContext, fileName)


					accountProxyDomainStatsMap[accountProxyStatId][domain] = domainProxyStats

					
				}
				


			}


			log.Println("accountProxyStatsMap", accountProxyStatsMap)
			log.Println("accountProxyDomainStatsMap", accountProxyDomainStatsMap)


			/*

				UPDATE DB

			*/


	
		}



		// delete activeProxyStatProcessingKey from redis --> allow other processes process it
		err = statsProxyRedisClient.Del(redisContext, activeProxyStatProcessingKey).Err() 
		if err != nil {
			errData := map[string]interface{}{
				"activeProxyStatProcessingKey": activeProxyStatProcessingKey,
			}
			logger.LogError("WARN", fileName, err, "failed to delete activeProxyStatProcessingKey in Redis", errData)
		}
		
	}


}

