package main

import (
	"go_proxy_worker/db"
	"go_proxy_worker/logger"
	"go_proxy_worker/slack"
	"go_proxy_worker/utils"
	"go_proxy_worker/workers"

	// "go_proxy_worker/redis"
	"os"
)

func main() {

	utils.LoadEnvVars()
	logger.UseJSONLogFormat()

	db.Init()
	db.InitCoreProxyRedis()
	db.InitStatsProxyRedis()
	db.InitConcurrencyProxyRedis()

	if utils.ProdEnv() {

		// Worker Startup Alert
		go slack.SlackStartupAlert()

	}

	args := os.Args[1:]
	if utils.ProdEnv() && len(args) == 0 {

		// PRODUCTION

		// NEW Redis
		// go workers.NEWCronSyncProxyConcurrency()
		// go slack.SlackTestAlert()
		// Workers
		go workers.CronCleanUserProxyConcurrency()
		go workers.CronUpdateAccountProxyStats()
		go workers.CronUpdateProxyStats()
		go workers.CronSyncProxyConcurrency()
		go workers.CronUpdatePPGBProxyStats()

		// Proxy Monitors
		go workers.CronCheckTopDomainPerformance()
		go workers.CronCheckProxyProviderFailedValidation()
		go workers.CronCheckEnterpriseUserPerformance()
		go workers.CronCheckProxyProviderDown()
		go workers.CronCheckProxyApiProfitability()
		go workers.CronCheckProxyProviderCredits()
		go workers.CronScrapeGithubRepo()
		go workers.CronScrapeYoutubeVideo()
		go workers.CronScrapeGoogleArticle()

	} else if !utils.ProdEnv() && len(args) == 0 {

		// In dev mode, use args to specify which worker to run.
		logger.LogTextSpace("DEVELOPMENT MODE: No args specified. Use args to specify which worker to run.")

	} else {
		if args[0] == "test" {

			// workers.RunTestScript()
			//workers.UpdateAccountProxyStats()
			workers.UpdateAccountProxyStats()

		}

		// DEVELOPMENT
		if args[0] == "dev" {

			workers.RunUpdatePPGBProxyStats()

			// workers.CronCleanUserProxyConcurrency()
			// workers.UpdateAccountProxyStats()
		}

		if args[0] == "testSlack" {
			slack.SlackTestAlert()
		}

		if args[0] == "accountStats" {
			workers.UpdateAccountProxyStats()
		}

		if args[0] == "proxySync" {
			workers.SyncProxyConcurrency()
		}

		if args[0] == "proxyAPICreditSync" {
			workers.SyncUserProxyAPICredits()
		}

		if args[0] == "proxyAPICreditSyncSINGLE" {
			workers.SyncSINGLEUserProxyAPICredits()
		}

		if args[0] == "checkUserProxyConcurrency" {
			workers.CheckUserProxyConcurrency()
		}

		if args[0] == "CheckProxyProviderCredits" {
			workers.CheckProxyProviderCredits()
		}

		if args[0] == "CleanUserProxyConcurrency" {
			workers.TestCleanUserProxyConcurrency()
		}

		if args[0] == "CheckProxyProviderFailedValidation" {
			workers.CheckProxyProviderFailedValidation()
		}

		if args[0] == "TestSlackFailedValidationAlert" {
			slack.TestSlackFailedValidationAlert()
		}

		if args[0] == "CheckEnterpriseUserPerformance" {
			workers.CheckEnterpriseUserPerformance()
		}

		if args[0] == "CheckTopDomainPerformance" {
			workers.CheckTopDomainPerformance()
		}

		if args[0] == "CheckProxyProviderDown" {
			workers.CheckProxyProviderDown()
		}

		if args[0] == "RunMonitors" {
			workers.CheckTopDomainPerformance()
			workers.CheckEnterpriseUserPerformance()
			workers.CheckProxyProviderFailedValidation()
			workers.CheckProxyProviderCredits()
		}

		if args[0] == "CheckProxyApiProfitability" {
			workers.CheckProxyApiProfitability()
		}

		if args[0] == "CheckUnpaidInvoices" {
			workers.CheckUnpaidInvoices()
		}

		if args[0] == "CheckFraudulentAccounts" {
			workers.CheckFraudulentAccounts()
		}

		if args[0] == "RunProxyTesterQueue" {
			workers.RunProxyTesterQueue()
		}

		if args[0] == "CheckGithubScraper" {
			workers.RunScrapeGithubRepo()
		}

		if args[0] == "CheckYoutubeScraper" {
			workers.RunScrapeYoutubeVideo()
		}

		if args[0] == "CheckArticleScraper" {
			workers.RunScrapeGoogleArticle()
		}
	}
}
