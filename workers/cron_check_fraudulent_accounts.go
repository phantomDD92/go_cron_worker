package workers

import (
	"fmt"
	"go_proxy_worker/db"
	"go_proxy_worker/logger"
	"os"
	"time"

	"go_proxy_worker/slack"

	"github.com/go-co-op/gocron"
)

type FraudluentAccount struct {
	Domain       string `json:"domain"`
	AccountIDs   string `json:"account_ids"`
	AccountNames string `json:"account_names"`
}

const BAN_REASON = "fraudulent_activity"

func CronCheckFraudulentAccounts() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Hours().Do(CheckFraudulentAccounts)
	s.StartBlocking()
}

func CheckFraudulentAccounts() {
	fileName := "cron_check_fraudulent_accounts.go"

	emptyErrMap := make(map[string]interface{})

	var db = db.GetDB()

	var faudulentAccounts []FraudluentAccount

	dbResult := db.Raw(`
			SELECT d.account_proxy_domain_stat_domain AS domain, 
       STRING_AGG(a.id::TEXT, ',') AS account_ids, 
       STRING_AGG(a.name, ',') AS account_names
			FROM account_proxy_domain_stats d
			JOIN accounts a ON d.account_id = a.id
			WHERE a.created_at >= CURRENT_DATE - INTERVAL '1 DAY'
				AND NOT EXISTS (
      		SELECT 1 
      		FROM banned_domains b 
      		WHERE b.domain = d.account_proxy_domain_stat_domain 
        	AND b.account_id = a.id
  			)
			GROUP BY d.account_proxy_domain_stat_domain
			HAVING COUNT(DISTINCT a.id) >= 3;
	`).Scan(&faudulentAccounts)
	if dbResult.Error != nil {
		logger.LogError("INFO", fileName, dbResult.Error, "failed to get FraudulentAccounts from DB", emptyErrMap)
	}

	EXPRESS_BACKEND_URL := os.Getenv("EXPRESS_BACKEND_URL")

	if len(faudulentAccounts) > 0 {
		for _, account := range faudulentAccounts {
			headline := fmt.Sprintf("Fraudulent Accounts Detected on %s", account.Domain)
			msg := fmt.Sprintf("%s fraudulent accounts scraping %s domain detected. Do you want to ban this domain for the newly created accounts?", account.AccountNames, account.Domain)
			approveLink := fmt.Sprintf("%s/v1/proxy/ban-domains?domain=%s&account_ids=%s&banned_reason=%s", EXPRESS_BACKEND_URL, account.Domain, account.AccountIDs, BAN_REASON)

			slack.SlackFraudAlert("#demo-notifications", headline, msg, approveLink)
		}
	}
}
