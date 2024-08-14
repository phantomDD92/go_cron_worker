package workers

import (
	"github.com/go-co-op/gocron"
	"go_proxy_worker/logger"
	"go_proxy_worker/utils"
	// "go_proxy_worker/dbRedisQueries"
	"github.com/go-redis/redis/v8"
	// "go_proxy_worker/slack"

	// "encoding/json"
	// "io/ioutil"
	"strconv"

	"go_proxy_worker/db"
	// "net/http"
	"time"
	"log"

)





func CronCleanUserProxyConcurrency() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(15).Seconds().Do(CleanUserProxyConcurrency)
	s.StartBlocking()
}




func CleanUserProxyConcurrency(){

	fileName := "cron_clean_user_proxy_concurrency.go"

	// emptyErrMap := make(map[string]interface{})

	// Redis Details
	var concurrencyProxyRedisClient = db.GetConcurrencyProxyRedisClient()
	// var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	redisAccountConcurrencyOverallListKey := "accountConcurrencyList?overall"
	listLength, err := concurrencyProxyRedisClient.SCard(redisContext, redisAccountConcurrencyOverallListKey).Result()
	if err != nil {
		errData := map[string]interface{}{
			"key": redisAccountConcurrencyOverallListKey,
		}
		logger.LogError("INFO", fileName, err, "redisAccountConcurrencyOverallListKey not in Redis", errData)
	}

	log.Println("listLength", listLength) 
	// log.Println("listLength", listLength) 
	for x := 0; x < int(listLength); x++ {

		accountKeyString, err := concurrencyProxyRedisClient.SPop(redisContext, redisAccountConcurrencyOverallListKey).Result()
		if err == nil {
			log.Println("accountKeyString", accountKeyString) 

			// time now
			now := time.Now()
			nowMinus3Minute := now.Add(-time.Minute * time.Duration(3)) 
			nowMinus3MinuteUnix := nowMinus3Minute.UnixNano()
			unixString := strconv.FormatInt(nowMinus3MinuteUnix, 10)

			concurrencyProxyRedisClient.ZRemRangeByScore(redisContext, accountKeyString, "-inf", unixString)
		}
	}
	
} 


func TestCleanUserProxyConcurrency2(){

	fileName := "cron_clean_user_proxy_concurrency.go"

	// emptyErrMap := make(map[string]interface{})

	// Redis Details
	var concurrencyProxyRedisClient = db.GetConcurrencyProxyRedisClient()
	// var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	redisAccountConcurrencyOverallListKey := "accountConcurrencyList?overall"
	listLength, err := concurrencyProxyRedisClient.SCard(redisContext, redisAccountConcurrencyOverallListKey).Result()
	if err != nil {
		errData := map[string]interface{}{
			"key": redisAccountConcurrencyOverallListKey,
		}
		logger.LogError("INFO", fileName, err, "redisAccountConcurrencyOverallListKey not in Redis", errData)
	}


	// log.Println("listLength", listLength) 
	for x := 0; x < int(listLength); x++ {

		accountKeyString, err := concurrencyProxyRedisClient.SPop(redisContext, redisAccountConcurrencyOverallListKey).Result()
		if err == nil {
			log.Println("accountKeyString", accountKeyString) 
		}
	}
	
} 


func TestCleanUserProxyConcurrency(){

	// fileName := "cron_clean_user_proxy_concurrency.go"

	// emptyErrMap := make(map[string]interface{})

	// Redis Details
	var concurrencyProxyRedisClient = db.GetConcurrencyProxyRedisClient()
	// var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	// time now
	now := time.Now()
	nowMinus3Minute := now.Add(-time.Minute * time.Duration(3)) 
	nowMinus3MinuteUnix := nowMinus3Minute.UnixNano()
	unixString := strconv.FormatInt(nowMinus3MinuteUnix, 10)

	accountKeyString := "accountConcurrencyList?account_id=9283"

	concurrencyProxyRedisClient.ZRemRangeByScore(redisContext, accountKeyString, "-inf", unixString)

	zRange := &redis.ZRangeBy{
		Min: "-inf",
		Max: unixString,
	}

	requests, err := concurrencyProxyRedisClient.ZRangeByScore(redisContext, accountKeyString, zRange).Result()
	if err != nil {
		log.Println("err", err)	
	}

	log.Println("requests", requests)	
	log.Println("len requests", len(requests))	
	log.Println("nowMinus3Minute", nowMinus3Minute)	
	log.Println("nowMinus3MinuteUnix", nowMinus3MinuteUnix)	
	
}
	



	

