package slack

import (
	"fmt"
	"go_proxy_worker/db"
	"go_proxy_worker/logger"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
	"github.com/stripe/stripe-go/v79"
)

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func SlackUnpaidInvoiceAlert(channel string, unpaidInvoices []*stripe.Invoice) {
	fileName := "slack_unpaid_invoice_alert.go"
	emptyErrMap := make(map[string]interface{})

	OAUTH_TOKEN := os.Getenv("SLACK_BEARER_TOKEN")
	api := slack.New(OAUTH_TOKEN)

	var db = db.GetDB()

	var accountsArray []Account
	accountResult := db.Raw(`
	SELECT id, name from accounts
	`).Scan(&accountsArray)
	if accountResult.Error != nil {
		logger.LogError("INFO", fileName, accountResult.Error, "failed to get accounts from DB", emptyErrMap)
	}
	accountsById := make(map[int]Account)
	for _, account := range accountsArray {
		accountsById[account.ID] = account
	}

	headline := "*The following enterprise invoices are currently outstanding:*"
	headerText := slack.NewTextBlockObject("mrkdwn", headline, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	invBlocks := make([]slack.Block, 0)
	invBlocks = append(invBlocks, headerSection)
	for _, inv := range unpaidInvoices {
		accountId, err := strconv.Atoi(inv.Metadata["account_id"])
		if err != nil {
			logger.LogError("INFO", fileName, err, "failed to convert account_id to int", emptyErrMap)
			continue
		}
		account := accountsById[accountId]
		created := time.Unix(inv.Created, 0)
		diff := time.Since(created)
		days := int(diff.Hours() / 24)
		invText := fmt.Sprintf("Account Name: %s\n Amount: $%d\n Days Outstanding: %d days\n", account.Name, inv.AmountDue/100, days)
		invField := slack.NewTextBlockObject("mrkdwn", invText, false, false)
		invBlocks = append(invBlocks, slack.NewSectionBlock(nil, []*slack.TextBlockObject{invField}, nil))
	}

	_, _, err := api.PostMessage(channel,
		slack.MsgOptionBlocks(invBlocks...))
	if err == nil {
		strErr := fmt.Sprintf("Message successfully sent to Channel %s at %s\n", channel, time.Now())
		logger.LogError("INFO", fileName, nil, strErr, emptyErrMap)
	} else {
		logger.LogError("ERROR", fileName, err, "Error Sending Slack Message", emptyErrMap)
	}
}
