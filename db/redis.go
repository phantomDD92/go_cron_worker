package db

import (
	"os"
	"log"
	// "time"
	"context"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)


var coreProxyRedisClient *redis.Client
var concurrencyProxyRedisClient *redis.Client
var statsProxyRedisClient *redis.Client

var (
	RedisCtx = context.TODO()
 )

func InitCoreProxyRedis() {


	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("CORE_PROXY_REDIS_CONNECTION_STRING")
	password := os.Getenv("CORE_PROXY_REDIS_PASSWORD")

	coreProxyRedisClient = redis.NewClient(&redis.Options{
		Addr: connectionString,
		Password: password,
		DB: 0,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			//Certificates: []tls.Certificate{cert}
		},
	})

	pong, err := coreProxyRedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Println("Core Proxy Redis Failed To Connect", err)
	} else {
		log.Println("Core Proxy Redis Connected:", pong)
	}

}

func InitConcurrencyProxyRedis() {


	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("CONCURRENCY_PROXY_REDIS_CONNECTION_STRING")
	password := os.Getenv("CONCURRENCY_PROXY_REDIS_PASSWORD")

	concurrencyProxyRedisClient = redis.NewClient(&redis.Options{
		Addr: connectionString,
		Password: password,
		DB: 0,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			//Certificates: []tls.Certificate{cert}
		},
	})

	pong, err := concurrencyProxyRedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Println("Concurrency Proxy Redis Failed To Connect", err)
	} else {
		log.Println("Concurrency Proxy Redis Connected:", pong)
	}

}

func InitStatsProxyRedis() {


	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("STATS_PROXY_REDIS_CONNECTION_STRING")
	password := os.Getenv("STATS_PROXY_REDIS_PASSWORD")

	statsProxyRedisClient = redis.NewClient(&redis.Options{
		Addr: connectionString,
		Password: password,
		DB: 0,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			//Certificates: []tls.Certificate{cert}
		},
	})

	pong, err := statsProxyRedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Println("Stats Proxy Redis Failed To Connect", err)
	} else {
		log.Println("Stats Proxy Redis Connected:", pong)
	}

}


func GetCoreProxyRedisClient() *redis.Client {
	return coreProxyRedisClient
}

func GetConcurrencyProxyRedisClient() *redis.Client {
	return concurrencyProxyRedisClient
}

func GetStatsProxyRedisClient() *redis.Client {
	return statsProxyRedisClient
}




