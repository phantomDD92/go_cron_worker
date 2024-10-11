package workers

import (
	"context"
	"fmt"
	"go_proxy_worker/db"
	"go_proxy_worker/dbRedisQueries"
	"go_proxy_worker/logger"
	"go_proxy_worker/models"
	"go_proxy_worker/utils"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"

	"time"
)

func CronUpdatePPGBProxyStats() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minutes().Do(RunUpdatePPGBProxyStats)
	s.StartBlocking()
}

func parseAccountID(s string) (string, error) {
	// Split the string into two parts: before and after the '?'
	parts := strings.Split(s, "?")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format")
	}

	// Parse the query part
	queryPart := parts[1]
	values, err := url.ParseQuery(queryPart)
	if err != nil {
		return "", err
	}

	// Extract the account_id
	accountIDs := values["account_id"]
	if len(accountIDs) == 0 {
		return "", fmt.Errorf("account_id not found")
	}

	return accountIDs[0], nil
}

// Helper function to get and parse the total requests/credits/bandwidth per day
func GetTotalsPerDay(redisClient *redis.Client, ctx context.Context, key string, fileName string) uint64 {
	totalsPerDateString, err := redisClient.Get(ctx, key).Result()

	if err != nil {
		errData := map[string]interface{}{
			"totalsPerDateString": key,
		}
		logger.LogError("INFO", fileName, err, "totalsPerDateString not in Redis", errData)
		return 0
	}

	totalsPerDateInt := utils.StringToUint64(totalsPerDateString)

	return totalsPerDateInt
}

func parseStats(input string) (int, int, int, int, int, error) {
	var numRequests, numSuccessfulRequests, numFailedRequests, creditsUsed, bandwidthUsed int

	// Splitting the string by "&"
	params := strings.Split(input, "&")
	for _, p := range params {
		// Splitting each parameter by "="
		kv := strings.Split(p, "=")
		if len(kv) != 2 {
			return 0, 0, 0, 0, 0, fmt.Errorf("invalid parameter: %s", p)
		}

		key, valueStr := kv[0], kv[1]
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return 0, 0, 0, 0, 0, fmt.Errorf("invalid value for %s: %s", key, valueStr)
		}

		// Assigning the value based on the key
		switch key {
		case "requests":
			numRequests = value
		case "successful_requests":
			numSuccessfulRequests = value
		case "failed_requests":
			numFailedRequests = value
		case "credits_used":
			creditsUsed = value
		case "bandwidth_used":
			bandwidthUsed = value
		}
	}

	return numRequests, numSuccessfulRequests, numFailedRequests, creditsUsed, bandwidthUsed, nil
}

// STEPS
// 1. Get all the active accounts which have processed a request in the last 24 hours (usually less than 1000)
// 2. Loop through these accounts and update the DB for
// - the total bytes used for that account - totalBandwidthUsed?api_key=*
// - the total bytes used for the current date - totalBaandwidthUsed?api_key=" + apiKey + "&date=" + dateString
// - the total bytes used for the current date & domain - totalBandwidthUsed?api_key=" + apiKey + "&date=" + dateString + "&domain=" + domain
func RunUpdatePPGBProxyStats() {

	fileName := "cron_update_sops_proxy_stats.go"

	// Redis Details
	var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	// load DB
	var db = db.GetDB()

	// Get Active Proxy Accounts which have send requests in the last day
	activePPGBAccountKeySet := "activeProxyPPGBKeySet"

	// listActivePPGBAccountsKeys, err := statsProxyRedisClient.SMembers(redisContext, activePPGBAccountKeySet).Result()
	// if err != nil {
	// 	errData := map[string]interface{}{
	// 		"activeProxyPPGBKeySet": activePPGBAccountKeySet,
	// 	}
	// 	logger.LogError("INFO", fileName, err, "activeProxyPPGBKeySet not in Redis", errData)
	// }

	var listActivePPGBAccountsKeys []string
	var err error

	runTestAccounts, listTestAccounts := utils.OnlyRunTestAccounts()
	if runTestAccounts {
		listActivePPGBAccountsKeys = listTestAccounts
		log.Println("listActivePPGBAccountsKeys TEST ACCOUNTS", listActivePPGBAccountsKeys)
	} else {
		listActivePPGBAccountsKeys, err = statsProxyRedisClient.SMembers(redisContext, activePPGBAccountKeySet).Result()
		if err != nil {
			errData := map[string]interface{}{
				"activePPGBAccountKeySet": activePPGBAccountKeySet,
			}
			logger.LogError("INFO", fileName, err, "activePPGBAccountKeySet not in Redis", errData)
		}
	}

	log.Println("************ listActivePPGBAccountsKeys **********")
	log.Println("listActivePPGBAccountsKeys", listActivePPGBAccountsKeys)
	log.Println("************ listActivePPGBAccountsKeys END **********")

	// loop through the accounts
	for _, accountId := range listActivePPGBAccountsKeys {

		//Get the account ppgb stats from redis
		if accountId != "" && (accountId != "3" || runTestAccounts) {
			accountTotalBandwidthUsedKey := "ppgbTotalBandwidthUsed?account_id=" + accountId
			accountTotalCreditsUsedKey := "ppgbTotalCreditsUsed?account_id=" + accountId

			// Account Details
			accountIdInt := utils.StringToUint(accountId)
			accountItem, validAccount := dbRedisQueries.GetAccountDetails(accountIdInt, db, coreProxyRedisClient, redisContext, fileName)

			if validAccount {
				/////////////
				// STEP 1 - Updating the total bandwidth / credits used in redis & the DB
				/////////////

				//get account date based on account timezone offset
				dateLayout := "2006-01-02"
				userDateTime, err := utils.GetCurrentDateTimeForOffset(accountItem)
				if err != nil {
					println("Failed to get the correct date offset:", err)
					// errData := structs.Map(sopsProxyProvider)
					// logger.LogError("ERROR", fileName, err, "Failed to update or create sopsProxyProvider in DB", errData)
				}
				usersDateString := userDateTime.Format(dateLayout)
				////////

				//accountTotalBandwidthUsedValue
				accountTotalBandwidthUsedValue, err := coreProxyRedisClient.Get(redisContext, accountTotalBandwidthUsedKey).Result()
				if err != nil {
					errData := map[string]interface{}{
						"accountTotalBandwidthUsedValue": accountTotalBandwidthUsedValue,
					}
					logger.LogError("INFO", fileName, err, "accountTotalBandwidthUsedValue not in Redis", errData)
				}

				//accountTotalCreditsUsedValue
				accountTotalCreditsUsedValue, err := coreProxyRedisClient.Get(redisContext, accountTotalCreditsUsedKey).Result()
				if err != nil {
					errData := map[string]interface{}{
						"accountTotalCreditsUsedValue": accountTotalCreditsUsedValue,
					}
					logger.LogError("INFO", fileName, err, "accountTotalCreditsUsedValue not in Redis", errData)
				}

				// Parse the bandwidth string to uint64
				accountTotalBandwidthInt := utils.StringToUint64(accountTotalBandwidthUsedValue)

				// Convert uint64 to uint (if it's safe to do so)
				accountTotalBandwidthUInt := uint(accountTotalBandwidthInt)

				accountIdInt := utils.StringToUint64(accountId)

				//Then try and get create/update the account_ppgb_overall_stats row -  if it does not exist then we need to create it
				var ppgbOverallStats models.AccountPPGBOverallStats
				ppgbOverallStatsRow := db.Table("account_ppgb_overall_stats").Where("account_id = ?", accountId).First(&ppgbOverallStats)
				if ppgbOverallStatsRow.Error != nil {

					println("Creating a new account_ppgb_overall_stats row...")
					// logger.LogError("ERROR", fileName, err, "Failed to get account_ppgb_overall_stats row in DB", accountId)

					ppgbOverallStats.AccountID = uint(accountIdInt)
					ppgbOverallStats.Date = userDateTime
					ppgbOverallStats.SuccessfulRequests = 1
					ppgbOverallStats.FailedRequests = 0
					ppgbOverallStats.TotalRequests = 1
					ppgbOverallStats.BytesUsed = accountTotalBandwidthUInt
					ppgbOverallStats.CreditsUsed = accountTotalBandwidthUInt
					// ppgbOverallStats.Type = "residential"

					createResult := db.Create(&ppgbOverallStats)
					if createResult.Error != nil {
						println("Failed to update or create account_ppgb_overall_stats in DB")
						// errData := structs.Map(sopsProxyProvider)
						// logger.LogError("ERROR", fileName, err, "Failed to update or create sopsProxyProvider in DB", errData)
					}

				} else {
					println("Updating account_ppgb_overall_stats row...")

					//if the row exists then we just update the value
					ppgbOverallStats.CreditsUsed = accountTotalBandwidthUInt
					ppgbOverallStats.BytesUsed = accountTotalBandwidthUInt
					ppgbOverallStats.TotalRequests = ppgbOverallStats.TotalRequests + 1
					updateResult := db.Save(&ppgbOverallStats)
					if updateResult.Error != nil || updateResult.RowsAffected == 0 {
						println("Failed to update account_ppgb_overall_stats in DB")
						// 		errData := structs.Map(dayProxyStat)
						// 		logger.LogError("ERROR", fileName, err, "Failed to update or create account_ppgb_overall_stats in DB", errData)

					}
				}

				/////////////
				// STEP 2 - Updating the DAILY total bandwidth used in redis & the DB
				/////////////

				//dailyTotalRequests
				redisTotalRequestsPerDateString := "ppgbTotalRequestsPerDate?account_id=" + accountId + "&date=" + usersDateString
				totalRequestsPerDateInt := GetTotalsPerDay(statsProxyRedisClient, redisContext, redisTotalRequestsPerDateString, fileName)

				//dailyTotalCredits
				redisTotalCreditsUsedPerDateString := "ppgbTotalCreditsUsedPerDate?account_id=" + accountId + "&date=" + usersDateString
				totalCreditsUsedPerDateInt := GetTotalsPerDay(statsProxyRedisClient, redisContext, redisTotalCreditsUsedPerDateString, fileName)

				//dailyTotalBytes
				redisTotalBandwidthUsedPerDateString := "ppgbtotalBandwidthUsedPerDate?account_id=" + accountId + "&date=" + usersDateString
				totalBandwidthUsedPerDateInt := GetTotalsPerDay(statsProxyRedisClient, redisContext, redisTotalBandwidthUsedPerDateString, fileName)

				//Then try and get create/update the account_ppgb_daily_stats row -  if it does not exist then we need to create it
				var ppgbDailyStats models.AccountPPGBDailyStats
				ppgbStatsDailyRow := db.Table("account_ppgb_daily_stats").Where("account_id = ? and date = ?", accountId, usersDateString).First(&ppgbDailyStats)
				if ppgbStatsDailyRow.Error != nil {

					println("Creating a new account_ppgb_daily_stats row...")
					// logger.LogError("ERROR", fileName, err, "Failed to get account_ppgb_overall_stats row in DB", accountId)
					// TEMP - TODO - FIX THIS TO USE THE ABOVE

					ppgbDailyStats.AccountID = uint(accountIdInt)
					ppgbDailyStats.Date = userDateTime
					ppgbDailyStats.SuccessfulRequests = uint(totalRequestsPerDateInt)
					ppgbDailyStats.FailedRequests = 0
					ppgbDailyStats.TotalRequests = uint(totalRequestsPerDateInt)
					ppgbDailyStats.BytesUsed = uint(totalBandwidthUsedPerDateInt)
					ppgbDailyStats.CreditsUsed = uint(totalCreditsUsedPerDateInt)
					// ppgbDailyStats.Type = "residential" //residential or mobile

					createResult := db.Create(&ppgbDailyStats)
					if createResult.Error != nil {
						println("Failed to update or create account_ppgb_daily_stats in DB")
						// errData := structs.Map(sopsProxyProvider)
						// logger.LogError("ERROR", fileName, err, "Failed to update or create sopsProxyProvider in DB", errData)
					}

				} else {

					//if the row exists then we just update the value
					ppgbDailyStats.TotalRequests = uint(totalRequestsPerDateInt)
					ppgbDailyStats.SuccessfulRequests = uint(totalRequestsPerDateInt)
					ppgbDailyStats.BytesUsed = uint(totalBandwidthUsedPerDateInt)
					ppgbDailyStats.CreditsUsed = uint(totalCreditsUsedPerDateInt)

					updateResult := db.Save(&ppgbDailyStats)
					if updateResult.Error != nil || updateResult.RowsAffected == 0 {
						println("Failed to update account_ppgb_overall_stats in DB")
						// 		errData := structs.Map(dayProxyStat)
						// 		logger.LogError("ERROR", fileName, err, "Failed to update or create account_ppgb_overall_stats in DB", errData)

					}
				}

				/////////////
				// STEP 3 - Updating the account DAILY DOMAIN stats
				/////////////

				// Get Active Proxy Accounts which have send requests in the last day
				activePPGBAccountDomainKeySet := "activePPGBAccountDomainKeySet?account_id=" + accountId
				listActivePPGBAccountDomainKeySet, err := statsProxyRedisClient.SMembers(redisContext, activePPGBAccountDomainKeySet).Result()
				if err != nil {
					errData := map[string]interface{}{
						"listActivePPGBAccountDomainKeySet": listActivePPGBAccountDomainKeySet,
					}
					logger.LogError("INFO", fileName, err, "listActivePPGBAccountDomainKeySet not in Redis", errData)
				}

				log.Println("listActivePPGBAccountDomainKeySet", listActivePPGBAccountDomainKeySet)

				// loop through the domain stats for this account & save them in the DB
				for _, activeDomain := range listActivePPGBAccountDomainKeySet {

					redisPPGBAccountDateAndDomainDetailsString := "ppgbOverallAccountPPGBProxyDetails?account_id=" + accountId + "&date=" + usersDateString + "&domain=" + activeDomain
					println("**** redis key ****")
					println(redisPPGBAccountDateAndDomainDetailsString)

					redisPPGBAccountDateAndDomainDetailsValues, err := statsProxyRedisClient.Get(redisContext, redisPPGBAccountDateAndDomainDetailsString).Result()
					if err != nil {
						errData := map[string]interface{}{
							"redisPPGBAccountDateAndDomainDetailsValues": redisPPGBAccountDateAndDomainDetailsValues,
						}
						logger.LogError("INFO", fileName, err, "redisPPGBAccountDateAndDomainDetailsString not in Redis", errData)

					}

					//parse the values
					var (
						numRequests           int
						numSuccessfulRequests int
						numFailedRequests     int
						creditsUsed           int
						bandwidthUsed         int
					)

					numRequests, numSuccessfulRequests, numFailedRequests, creditsUsed, bandwidthUsed, err = parseStats(redisPPGBAccountDateAndDomainDetailsValues)

					//Then try and get create/update the row -  if it does not exist then we need to create it
					var ppgbDailyDomainStats models.AccountPPGBDailyDomainStats
					ppgbStatsDailyRow := db.Table("account_ppgb_daily_domain_stats").Where("account_id = ? and date = ? and domain = ?", accountId, usersDateString, activeDomain).First(&ppgbDailyDomainStats)
					if ppgbStatsDailyRow.Error != nil {

						println("Creating a new account_ppgb_daily_domain_stats row...")
						// logger.LogError("ERROR", fileName, err, "Failed to get account_ppgb_overall_stats row in DB", accountId)

						ppgbDailyDomainStats.AccountID = uint(accountIdInt)
						ppgbDailyDomainStats.Date = userDateTime
						ppgbDailyDomainStats.SuccessfulRequests = uint(numSuccessfulRequests)
						ppgbDailyDomainStats.FailedRequests = uint(numFailedRequests)
						ppgbDailyDomainStats.TotalRequests = uint(numRequests)
						ppgbDailyDomainStats.BytesUsed = uint(bandwidthUsed)
						ppgbDailyDomainStats.CreditsUsed = uint(creditsUsed)
						ppgbDailyDomainStats.Domain = activeDomain
						ppgbDailyDomainStats.AccountPPGBDailyStatID = ppgbDailyStats.ID

						createResult := db.Create(&ppgbDailyDomainStats)
						if createResult.Error != nil {
							println("Failed to update or create account_ppgb_daily_stats in DB")
							// errData := structs.Map(sopsProxyProvider)
							// logger.LogError("ERROR", fileName, err, "Failed to update or create sopsProxyProvider in DB", errData)
						}

					} else {
						println("Updating account_ppgb_daily_domain_stats row...")

						//if the row exists then we just update the value
						ppgbDailyDomainStats.SuccessfulRequests = uint(numSuccessfulRequests)
						ppgbDailyDomainStats.FailedRequests = uint(numFailedRequests)
						ppgbDailyDomainStats.TotalRequests = uint(numRequests)
						ppgbDailyDomainStats.BytesUsed = uint(bandwidthUsed)
						ppgbDailyDomainStats.CreditsUsed = uint(bandwidthUsed)

						updateResult := db.Save(&ppgbDailyDomainStats)
						if updateResult.Error != nil || updateResult.RowsAffected == 0 {
							println("Failed to update account_ppgb_daily_domain_stats in DB")
							// 		errData := structs.Map(dayProxyStat)
							// 		logger.LogError("ERROR", fileName, err, "Failed to update or create account_ppgb_overall_stats in DB", errData)

						}
					}

					// Remove From listActivePPGBAccountDomainKeySet
					err = statsProxyRedisClient.SRem(redisContext, activePPGBAccountDomainKeySet, activeDomain).Err()
					if err != nil {
						errData := map[string]interface{}{
							"apiKey":                            activeDomain,
							"listActivePPGBAccountDomainKeySet": activePPGBAccountDomainKeySet,
						}
						logger.LogError("ERROR", fileName, err, "failed to delete activeDomain from activePPGBAccountDomainKeySet in Redis", errData)
					}
				}

				// Remove accountId From activeAccountAccountsKeySet
				err = statsProxyRedisClient.SRem(redisContext, activePPGBAccountKeySet, accountId).Err()
				if err != nil {
					errData := map[string]interface{}{
						"accountId":             accountId,
						"activeProxyPPGBKeySet": activePPGBAccountKeySet,
					}
					logger.LogError("ERROR", fileName, err, "failed to delete accountId from activeProxyPPGBKeySet in Redis", errData)
				}

			}

		}

	}

}
