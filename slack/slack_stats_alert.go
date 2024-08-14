package slack

import (
	"fmt"
	"go_proxy_worker/logger"
	"os"

	"github.com/slack-go/slack"
)

func TestSlackFailedValidationAlert() {

	headline := "TEST"

	fileName := "slack_failed_validation_alert"
	emptyErrMap := make(map[string]interface{})

	OAUTH_TOKEN := os.Getenv("SLACK_BEARER_TOKEN")
	CHANNEL_ID := "#proxy-provider-failed-validation"

	// Header Section
	var headerText *slack.TextBlockObject
	headerString := "*" + headline + "*"
	headerText = slack.NewTextBlockObject("mrkdwn", headerString, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	statsString := "```Domain | Req | Success | Invalid | % |\n"
	statsString = statsString + "test.com | 1230 | 1230 | 1230 | 10% |\n```"
	// statsString = statsString + fmt.Sprintf("%v", jobStruct.SpiderName) + "```"
	statsField := slack.NewTextBlockObject("mrkdwn", statsString, false, false)

	fieldSlice1 := make([]*slack.TextBlockObject, 0)
	fieldSlice1 = append(fieldSlice1, statsField)
	rowSection1 := slack.NewSectionBlock(nil, fieldSlice1, nil)

	api := slack.New(OAUTH_TOKEN)

	previewMessage := headline

	channelId, timestamp, err := api.PostMessage(
		CHANNEL_ID,
		slack.MsgOptionText(previewMessage, false),
		slack.MsgOptionBlocks(
			headerSection,
			rowSection1,
		),
		slack.MsgOptionAsUser(true),
	)

	if err == nil {
		strErr := fmt.Sprintf("Message successfully sent to Channel %s at %s\n", channelId, timestamp)
		logger.LogError("INFO", fileName, nil, strErr, emptyErrMap)
	} else {
		logger.LogError("ERROR", fileName, err, "Error Sending Slack Message", emptyErrMap)
	}

}

func SlackStatsAlert(channel string, headline string, statsString string) {

	fileName := "slack_stats_alert"
	emptyErrMap := make(map[string]interface{})

	OAUTH_TOKEN := os.Getenv("SLACK_BEARER_TOKEN")
	CHANNEL_ID := channel

	// Header Section
	var headerText *slack.TextBlockObject
	headerString := "*" + headline + "*"
	headerText = slack.NewTextBlockObject("mrkdwn", headerString, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Stats Block
	statsField := slack.NewTextBlockObject("mrkdwn", statsString, false, false)
	fieldSlice1 := make([]*slack.TextBlockObject, 0)
	fieldSlice1 = append(fieldSlice1, statsField)
	rowSection1 := slack.NewSectionBlock(nil, fieldSlice1, nil)

	api := slack.New(OAUTH_TOKEN)

	previewMessage := headline

	channelId, timestamp, err := api.PostMessage(
		CHANNEL_ID,
		slack.MsgOptionText(previewMessage, false),
		slack.MsgOptionBlocks(
			headerSection,
			rowSection1,
		),
		slack.MsgOptionAsUser(true),
	)

	if err == nil {
		strErr := fmt.Sprintf("Message successfully sent to Channel %s at %s\n", channelId, timestamp)
		logger.LogError("INFO", fileName, nil, strErr, emptyErrMap)
	} else {
		logger.LogError("ERROR", fileName, err, "Error Sending Slack Message", emptyErrMap)
	}

}
