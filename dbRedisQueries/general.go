package dbRedisQueries

import (
	"strconv"
	"context"
	"go_proxy_worker/logger"
	"github.com/go-redis/redis/v8"
)


func GetUintRedis(key string, redisClient *redis.Client, redisCtx context.Context, fileName string) uint {

	redisString, _ := redisClient.Get(redisCtx, key).Result()
	// if err != nil {
	// 	errData := map[string]interface{}{
	// 		"key": key,
	// 	}
	// 	logger.LogError("INFO", fileName, err, (key + " not in Redis"), errData)
	// }

	// If No Value
	if redisString == "" {
		return 0
	}

	valUint64, err := strconv.ParseUint(redisString, 10, 32)
	if err != nil {
		errData := map[string]interface{}{}
		logger.LogError("INFO", fileName, err, "Error parsing string to Uint", errData)
	}

	return uint(valUint64)

}


func GetInt64Redis(key string, redisClient *redis.Client, redisCtx context.Context, fileName string) int64 {

	redisString, _ := redisClient.Get(redisCtx, key).Result()
	// if err != nil {
	// 	errData := map[string]interface{}{
	// 		"key": key,
	// 	}
	// 	logger.LogError("INFO", fileName, err, (key + " not in Redis"), errData)
	// }

	// If No Value
	if redisString == "" {
		return 0
	}

	i64, err := strconv.ParseInt(redisString, 10, 64)
	if err != nil {
		errData := map[string]interface{}{}
		logger.LogError("INFO", fileName, err, "Error parsing string to Uint", errData)
	}

	return i64

}


func GetFloat64Redis(key string, redisClient *redis.Client, redisCtx context.Context, fileName string) float64 {

	redisString, _ := redisClient.Get(redisCtx, key).Result()
	// if err != nil {
	// 	errData := map[string]interface{}{
	// 		"key": key,
	// 	}
	// 	logger.LogError("INFO", fileName, err, (key + " not in Redis"), errData)
	// }

	// If No Value
	if redisString == "" {
		return 0
	}

	valFloat64, err := strconv.ParseFloat(redisString, 32)
	if err != nil {
		errData := map[string]interface{}{}
		logger.LogError("INFO", fileName, err, "Error parsing string to Float64", errData)
	}

	return valFloat64

}


