package redis

import (
	"os"
	"log"
	"context"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)


var redisClient *redis.Client
var concurrencyProxyRedisClient *redis.Client
var err error

var (
	Ctx = context.TODO()
 )

func InitRedis() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("NEW_REDIS_CONNECTION_STRING")
	password := os.Getenv("NEW_REDIS_PASSWORD")

	redisClient = redis.NewClient(&redis.Options{
		Addr: connectionString,
		Password: password,
		DB: 0,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			//Certificates: []tls.Certificate{cert}
		},
	})

	pong, err := redisClient.Ping(Ctx).Result()
	if err != nil {
		log.Println("Redis failed to connect", err)
	} else {
		log.Println("Redis connected:", pong)
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

	pong, err := concurrencyProxyRedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Println("Concurrency Proxy Redis Failed To Connect", err)
	} else {
		log.Println("Concurrency Proxy Redis Connected:", pong)
	}

}


func GetRedisClient() *redis.Client {
	return redisClient
}

func GetConcurrencyProxyRedisClient() *redis.Client {
	return concurrencyProxyRedisClient
}
