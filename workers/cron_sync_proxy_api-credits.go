package workers

import (
	"github.com/go-co-op/gocron"
	"go_proxy_worker/logger"
	"go_proxy_worker/utils"
	"go_proxy_worker/models"
	"go_proxy_worker/db"
	"time"
	"fmt"
	"log"

)





func CronSyncUserProxyAPICredits() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Seconds().Do(SyncUserProxyAPICredits)
	s.StartBlocking()
}

type ProxyAccountDetails struct {
	AccountId					uint `json:"account_id"`
	StripeCustomerId			string `json:"stripe_customer_id"`
	StripeProductId				string `json:"stripe_product_id"`
	StripePriceId				string `json:"stripe_price_id"`
}

// type StripeSubscriptionData struct {
// 	Data	struct {
// 				Plan struct {
// 					Product string `json:"product"`
// 				} `json:"plan"`
// 				CurrentPeriodStart int `json:"current_period_start"`
// 			} `json:"data"`
// }

// type StripeSubscription struct {
// 	Data				[]StripeSubscriptionData `json:"data"`
// }

// type StripeResponse struct {
// 	Subscriptions				StripeSubscription `json:"subscriptions"`
// }



func SyncUserProxyAPICredits(){

	fileName := "cron_update_sops_proxy_stats.go"

	emptyErrMap := make(map[string]interface{})

	// Redis Details
	var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	// var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	// load DB
	var db = db.GetDB()

	// Get AccountProxy Rows From DB
	var accountProxyArray []models.AccountProxy
	accountProxyArrayResult := db.Table("account_proxy").Where("account_proxy_total_requests > ?", 0).Find(&accountProxyArray)
	if accountProxyArrayResult.Error != nil {
		logger.LogError("INFO", fileName, accountProxyArrayResult.Error, "failed to get account_proxy rows from DB", emptyErrMap)
	}

	for _, accountProxy := range accountProxyArray {

		// Update Redis
		redisProxyApiCreditsKey := "accountUsedApiCredits?account_id=" + fmt.Sprintf("%v", accountProxy.AccountId)
		err := coreProxyRedisClient.Set(redisContext, redisProxyApiCreditsKey, accountProxy.AccountProxyUsedCredits, 86400*time.Second).Err()
		if err != nil {
			logger.LogError("WARN", fileName, err, "failed to update accountUsedApiCredits in Redis", emptyErrMap)
		}

	}

}


func SyncSINGLEUserProxyAPICredits(){

	fileName := "cron_update_sops_proxy_stats.go"

	emptyErrMap := make(map[string]interface{})

	// Redis Details
	var coreProxyRedisClient = db.GetCoreProxyRedisClient()
	// var statsProxyRedisClient = db.GetStatsProxyRedisClient()
	redisContext := utils.GetRedisCtx()

	// load DB
	var db = db.GetDB()

	accountId := 4765

	// Get AccountProxy Rows From DB
	var accountProxyArray []models.AccountProxy
	accountProxyArrayResult := db.Table("account_proxy").Where("account_id = ?", accountId).Find(&accountProxyArray)
	if accountProxyArrayResult.Error != nil {
		logger.LogError("INFO", fileName, accountProxyArrayResult.Error, "failed to get account_proxy rows from DB", emptyErrMap)
	}

	

	for _, accountProxy := range accountProxyArray {

		log.Println("accountProxy.AccountProxyUsedCredits", accountProxy.AccountProxyUsedCredits)

		// Update Redis
		redisProxyApiCreditsKey := "accountUsedApiCredits?account_id=" + fmt.Sprintf("%v", accountProxy.AccountId)
		err := coreProxyRedisClient.Set(redisContext, redisProxyApiCreditsKey, accountProxy.AccountProxyUsedCredits, 86400*time.Second).Err()
		if err != nil {
			logger.LogError("WARN", fileName, err, "failed to update accountUsedApiCredits in Redis", emptyErrMap)
		}

	}

}


// func SyncUserProxyAPICredits(){

// 	/*
// 		GET ALL ACCOUNTS
// 	*/

// 	// Recalculate Job Groups From Jobs Data
// 	var proxyAccountDetailsArray []ProxyAccountDetails
// 	proxyAccountsResult := db.Raw(`
// 	SELECT 
// 	accounts.id as account_id,
// 	accounts.stripe_customer_id as stripe_customer_id,
// 	proxy_plans.stripe_product_id as stripe_product_id,
// 	proxy_plans.stripe_price_id as stripe_price_id
// 	FROM account_proxy
// 	JOIN accounts 
// 	ON account_proxy.account_id = accounts.id
// 	JOIN proxy_plans
// 	ON proxy_plans.id = accounts.proxy_plan_id
// 	WHERE account_proxy.account_proxy_total_requests > 0
// 	`, nil).Scan(&proxyAccountDetailsArray)

// 	for _, proxyAccountDetail := range proxyAccountDetailsArray {

// 		var currentPeriodStart int 

// 		if proxyAccountDetail.StripeCustomerId != "" {

// 			// Get Plan Start Date
// 			stripe.Key = "sk_test_51KsnoxBcAgNJKnIKAdaHDsRsAJT7Dok2sazz1zPPwbwY6AbE7tSC6BtWRCXNuDnMc6bSd3QFB4OYZtVHlGz50fnu00ZtH4JkV5"

// 			params := &stripe.ChargeParams{}
// 			params.AddExpand("subscriptions")
// 			customer, err := charge.Get(proxyAccountDetail.StripeCustomerId, params)

// 			subscriptionDataArray := customer.Subscriptions.Data 
// 			for _, subscription := range subscriptionData {

// 				if subscription.Plan.Product == proxyAccountDetail.StripeProductId {
// 					currentPeriodStart = subscription.CurrentPeriodStart
// 				}
// 			}

// 		}

// 	}

// 	for account in results[:1]:
//     account_id = account[0]
//     stripe_customer_id = account[1]
//     stripe_product_id = account[2]
//     stripe_plan_id = account[3]
//     proxy_subscription = {} 
//     start_date_epoch = None
    
//     if stripe_customer_id != '' and stripe_customer_id is not None:
        
//         ## Get Renewal Date
//         customer_object = stripe.Customer.retrieve(stripe_customer_id, expand=['subscriptions'])
//         subscription_data_array = customer_object.subscriptions.data
//         for subscription in subscription_data_array:
//             if subscription.plan.product == stripe_product_id:
//                 proxy_subscription = subscription
//                 start_date_epoch = proxy_subscription.current_period_start
                
        
//         ## Calculate
//         if renewal_date_epoch is not None:
//             dt_utc_aware = datetime.datetime.fromtimestamp(start_date_epoch, datetime.timezone.utc)
//             start_date = dt_utc_aware.strftime('%Y-%m-%d')
//             stats_calculation = get_all_query(f"""
//             SELECT 
//             SUM(account_proxy_stat_requests) as requests,
//             SUM(account_proxy_stat_successful) as successful,
//             SUM(account_proxy_stat_failed) as failed,
//             SUM(account_proxy_stat_credits)  as credits
//             from account_proxy_stats 
//             where account_id = {account_id} and account_proxy_stat_day_start_time >= '{start_date}' 
//             """)
            
        
//         if len(stats_calculation) > 0 and stats_calculation[0][3] is not None:
            
//             ## Update DB
//             update_API_credits = update_query(f"""
//             UPDATE account_proxy
//             SET account_proxy_used_credits = %s,
//             WHERE account_id = {account_id}
//             """,
//                 [
//                     stats_calculation[0][3],
//                 ])
        
//             ## Update Redis
//             api_credits_redis_key = 'accountUsedApiCredits?account_id=' + account_id

// }