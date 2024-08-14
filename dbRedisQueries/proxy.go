package dbRedisQueries

import (
	// "log"
	"fmt"
	"time"
	"context"
	"go_proxy_worker/utils"
	"go_proxy_worker/logger"
	"go_proxy_worker/models"
	"gorm.io/gorm"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/fatih/structs"
)



func GetAccountProxyStatId(accountId uint, usersDayStartTime time.Time, usersDayStartTimeString string, db *gorm.DB, redisClient *redis.Client, redisCtx context.Context, fileName string) uint {

	errData := map[string]interface{}{
		"accountId": accountId,
		"dayString": usersDayStartTimeString,
	}

	accountProxyStatMappingKey := "accountProxyStatMapping?account_id=" + fmt.Sprintf("%v", accountId) + "&dayStartDate=" + usersDayStartTimeString
	accountProxyStatMappingKey = utils.RedisEnvironmentVersion(accountProxyStatMappingKey)

	accountProxyStatMappingString, err := redisClient.Get(redisCtx, accountProxyStatMappingKey).Result()
	if err != nil {
		errData := map[string]interface{}{
			"accountProxyStatMappingKey": accountProxyStatMappingKey,
		}
		logger.LogError("INFO", fileName, err, "accountProxyStatMappingKey not in Redis", errData)
	}

	if accountProxyStatMappingString == "" {

		var accountProxyStat models.AccountProxyStat
		accountProxyStatResult := db.Table("account_proxy_stats").Where("account_id = ? and account_proxy_stat_day_start_time = ?", accountId, usersDayStartTime).First(&accountProxyStat)
		if accountProxyStatResult.Error != nil || accountProxyStatResult.RowsAffected == 0 {
			accountProxyStat.AccountId = accountId
			accountProxyStat.AccountProxyStatDayStartTime = usersDayStartTime
			createResult := db.Create(&accountProxyStat)
			if createResult.Error != nil {
				errData := structs.Map(accountProxyStat)
				logger.LogError("ERROR", fileName, err, "Failed to update or create accountProxyStats in DB", errData)
			}
		}

		// Storing Structs as JSON Strings
		idJSON, err := json.Marshal(accountProxyStat.ID)
		if err != nil {
			logger.LogError("WARN", fileName, err, "failed to parse accountProxyStat.Id Redis string to JSON", errData)
		}

		// update redis with Spider Stat
		err = redisClient.Set(redisCtx, accountProxyStatMappingKey, idJSON, 86400*time.Second).Err()
		if err != nil {
			logger.LogError("WARN", fileName, err, "failed to update accountProxyStat.Id in Redis", errData)
		}

		return accountProxyStat.ID

	}

	return utils.StringToUint(accountProxyStatMappingString)

	
}



func GetAccountProxy(accountId uint, db *gorm.DB, redisClient *redis.Client, redisCtx context.Context, fileName string) models.AccountProxy {

	errData := map[string]interface{}{
		"accountId": accountId,
	}

	accountProxyKey := "accountProxy?account_id=" + fmt.Sprintf("%v", accountId)

	accountProxyString, err := redisClient.Get(redisCtx, accountProxyKey).Result()
	if err != nil {
		errData := map[string]interface{}{
			"accountProxyKey": accountProxyKey,
		}
		logger.LogError("INFO", fileName, err, "accountProxyKey not in Redis", errData)
	}

	if accountProxyString == "" {

		var accountProxy models.AccountProxy
		accountProxyResult := db.Table("account_proxy").Where("account_id = ?", accountId).First(&accountProxy)
		if accountProxyResult.Error != nil {
			logger.LogError("ERROR", fileName, err, "Failed to get accountProxy in DB", errData)

		} else {

			// Storing Structs as JSON Strings
			accountProxyJSON, err := json.Marshal(accountProxy)
			if err != nil {
				logger.LogError("WARN", fileName, err, "failed to parse accountProxy Redis string to JSON", errData)
			}

			// update redis with Spider Stat
			err = redisClient.Set(redisCtx, accountProxyKey, accountProxyJSON, 86400*time.Second).Err()
			if err != nil {
				logger.LogError("WARN", fileName, err, "failed to update accountProxy in Redis", errData)
			}
		}

		

		return accountProxy

	} else {

		// If Redis Contains 
		var accountProxy models.AccountProxy
		json.Unmarshal([]byte(accountProxyString), &accountProxy)
		return accountProxy

	}

	return models.AccountProxy{}

	
}


func GetSopsProxyProviderId(proxy string, db *gorm.DB, redisClient *redis.Client, redisCtx context.Context, fileName string) uint {

	errData := map[string]interface{}{
		"proxy": proxy,
	}

	proxyProviderKey := "sopsProxyProviderMapping?proxyName=" + proxy

	proxyProviderString, err := redisClient.Get(redisCtx, proxyProviderKey).Result()
	if err != nil {
		errData := map[string]interface{}{
			"proxyProviderKey": proxyProviderKey,
		}
		logger.LogError("INFO", fileName, err, "proxyProviderKey not in Redis", errData)
	}

	if proxyProviderString == "" {

		var sopsProxyProvider models.SopsProxyProvider
		sopsProxyProviderResult := db.Table("sops_proxy_providers").Where("sops_proxy_name = ?", proxy).First(&sopsProxyProvider)
		if sopsProxyProviderResult.Error != nil {
			logger.LogError("ERROR", fileName, err, "Failed to get sopsProxyProvider in DB", errData)
			sopsProxyProvider.SopsProxyName = proxy
			createResult := db.Create(&sopsProxyProvider)
			if createResult.Error != nil {
				errData := structs.Map(sopsProxyProvider)
				logger.LogError("ERROR", fileName, err, "Failed to update or create sopsProxyProvider in DB", errData)
			}

		} else {

			// // Storing Structs as JSON Strings
			// accountProxyJSON, err := json.Marshal(sopsProxyProvider)
			// if err != nil {
			// 	logger.LogError("WARN", fileName, err, "failed to parse sopsProxyProvider Redis string to JSON", errData)
			// }

			// update redis with Spider Stat
			err = redisClient.Set(redisCtx, proxyProviderKey, sopsProxyProvider.ID, 86400*time.Second).Err()
			if err != nil {
				logger.LogError("WARN", fileName, err, "failed to update sopsProxyProvider in Redis", errData)
			}
		}

		

		return sopsProxyProvider.ID

	} else {

		return utils.StringToUint(proxyProviderString)

	}

	return 0

}


func GetSopsDayProxyStatId(sopsProxyProviderId uint, proxy string, dayStartTime time.Time, dayStartTimeString string, db *gorm.DB, redisClient *redis.Client, redisCtx context.Context, fileName string) uint {

	errData := map[string]interface{}{
		"sopsProxyProviderId": sopsProxyProviderId,
		"proxy": proxy,
		"dayString": dayStartTimeString,
	}

	sopsDayProxyStatMapping := "sopsDayProxyStatMapping?proxyProviderId=" + fmt.Sprintf("%v", sopsProxyProviderId) + "&dayStartDate=" + dayStartTimeString
	sopsDayProxyStatMapping = utils.RedisEnvironmentVersion(sopsDayProxyStatMapping)

	sopsDayProxyStatMappingString, err := redisClient.Get(redisCtx, sopsDayProxyStatMapping).Result()
	if err != nil {
		errData := map[string]interface{}{
			"sopsDayProxyStatMapping": sopsDayProxyStatMapping,
		}
		logger.LogError("INFO", fileName, err, "sopsDayProxyStatMapping not in Redis", errData)
	}


	if sopsDayProxyStatMappingString == "" {

		var sopsDayProxyStat models.SopsDayProxyStat
		sopsDayProxyStatResult := db.Table("sops_day_proxy_stats").Where("sops_proxy_provider_id = ? and sops_day_proxy_stat_day_start_time = ?", sopsProxyProviderId, dayStartTime).First(&sopsDayProxyStat)
		if sopsDayProxyStatResult.Error != nil || sopsDayProxyStatResult.RowsAffected == 0 {
			sopsDayProxyStat.SopsProxyProviderId = sopsProxyProviderId
			sopsDayProxyStat.SopsProxyName = proxy
			sopsDayProxyStat.SopsDayProxyStatDayStartTime = dayStartTime
			createResult := db.Create(&sopsDayProxyStat)
			if createResult.Error != nil {
				errData := structs.Map(sopsDayProxyStat)
				logger.LogError("ERROR", fileName, err, "Failed to update or create sopsDayProxyStat in DB", errData)
			}
		}

		// update redis with SopsDayProxyStat
		err = redisClient.Set(redisCtx, sopsDayProxyStatMapping, sopsDayProxyStat.ID, 86400*time.Second).Err()
		if err != nil {
			logger.LogError("WARN", fileName, err, "failed to update sopsDayProxyStat.ID in Redis", errData)
		}

		return sopsDayProxyStat.ID

	}

	return utils.StringToUint(sopsDayProxyStatMappingString)

}






