package workers

import (
	"github.com/go-co-op/gocron"
	"go_proxy_worker/logger"
	"go_proxy_worker/utils"
	"go_proxy_worker/db"
	"strconv"
	"time"
	"fmt"
	// "log"

)





func CronCheckUserProxyConcurrency() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Seconds().Do(CheckUserProxyConcurrency)
	s.StartBlocking()
}




func CheckUserProxyConcurrency(){

	fileName := "cron_check_user_proxy_concurrency.go"

	// emptyErrMap := make(map[string]interface{})

	// Redis Details
	var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	// var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()


	// List of User Accounts
	accountIdList := []uint{
		971,
		1221,
	}

	for accountId := range accountIdList {

		redisAccountConcurrencyKey := "accountConcurrency?account_id=" + fmt.Sprintf("%v", accountId)

		// Check Redis
		redisAccountConcurrencyBytes, err := coreProxyRedisClient.Get(redisContext, redisAccountConcurrencyKey).Result()
		if err != nil {
			errData := map[string]interface{}{
				"accountId": accountId,
			}
			logger.LogError("WARN", fileName, err, "accountId not in Redis", errData)
		}


		if len(redisAccountConcurrencyBytes) > 0 {

			intVar, err := strconv.Atoi(redisAccountConcurrencyBytes)
			if err != nil {
				errData := map[string]interface{}{
					"accountId": accountId,
				}
				logger.LogError("ERROR", fileName, err, "error converting redisAccountConcurrencyBytes to int", errData)
			}

			if err == nil && intVar < 0 {
				err = coreProxyRedisClient.Set(redisContext, redisAccountConcurrencyKey, 0, 5*60*time.Second).Err()
				if err != nil {
					errData := map[string]interface{}{
						"accountId": accountId,
					}
					logger.LogError("ERROR", fileName, err, "failed to set user proxy concurrency to 0 in Redis", errData)
				}
			}
			
		}


		
	}





}