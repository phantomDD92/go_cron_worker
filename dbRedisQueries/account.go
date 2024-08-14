package dbRedisQueries

import (
	"fmt"
	//"log"
	"time"
	"context"
	"go_proxy_worker/logger"
	"go_proxy_worker/models"
	"gorm.io/gorm"
	"github.com/go-redis/redis/v8"
	"encoding/json"
)


func GetAccountDetails(accountID uint, db *gorm.DB, redisClient *redis.Client, redisCtx context.Context, fileName string) (models.Account, bool) {

	emptyErrMap := make(map[string]interface{})


	/*
		CHECK REDIS, THEN POSTGRES
	*/

	
	redisAccountIdString := "account?id=" + fmt.Sprintf("%v", accountID)


	// Check Redis
	redisAccountString, err := redisClient.Get(redisCtx, redisAccountIdString).Result()
	if err != nil {
		errData := map[string]interface{}{
			"accountID": accountID,
		}
		logger.LogError("WARN", fileName, err, "accountID not in Redis", errData)
	}


	if redisAccountString == "" {

		// get account info
		var accountItem models.Account
		accountResult := db.Where("id = ?", accountID).First(&accountItem)
		if accountResult.Error == nil {
			
			// Storing Structs as JSON Strings
			accountItemJSON, err := json.Marshal(accountItem)
			if err != nil {
				logger.LogError("WARN", fileName, err, "failed to parse AccountItem Redis string to JSON", emptyErrMap)
			}

			// update redis with account item
			err = redisClient.Set(redisCtx, redisAccountIdString, accountItemJSON, 86400*time.Second).Err()
			if err != nil {
				logger.LogError("WARN", fileName, err, "failed to update AccountItem in Redis", emptyErrMap)
			}

			return accountItem, true

		} else if accountResult.Error != nil {

			errData := map[string]interface{}{
				"accountID": accountID,
			}
			logger.LogError("WARN", fileName, accountResult.Error, "issue getting account with ID from DB", errData)
			return accountItem, false

		} else if accountResult.RowsAffected == 0 {

			errData := map[string]interface{}{
				"accountID": accountID,
			}
			logger.LogError("WARN", fileName, accountResult.Error, "sent account ID not in DB.", errData)
			return accountItem, false

		}

	} else {

		// If Redis Contains 
		var accountItem models.Account
		json.Unmarshal([]byte(redisAccountString), &accountItem)
		return accountItem, true
	}


	// Else
	var emptyAccountItem models.Account
	return emptyAccountItem, false
	
	
}

