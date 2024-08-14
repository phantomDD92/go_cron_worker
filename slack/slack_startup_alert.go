package slack

import (
	"fmt"
	"go_proxy_worker/logger"
	"os"

	"github.com/slack-go/slack"
)

func SlackStartupAlert() {

	fileName := "slack_error_logger"
	emptyErrMap := make(map[string]interface{})

	OAUTH_TOKEN := os.Getenv("SLACK_BEARER_TOKEN")
	CHANNEL_ID := "#go-worker-monitor"

	// Header Section
	var headerText *slack.TextBlockObject
	headerString := "*GO PROXY WORKER - STARTED*"
	headerText = slack.NewTextBlockObject("mrkdwn", headerString, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	api := slack.New(OAUTH_TOKEN)

	previewMessage := "*GO PROXY WORKER - STARTED*"

	channelId, timestamp, err := api.PostMessage(
		CHANNEL_ID,
		slack.MsgOptionText(previewMessage, false),
		slack.MsgOptionBlocks(
			headerSection,
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
