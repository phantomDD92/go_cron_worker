package slack

import (
	"fmt"
	"go_proxy_worker/logger"
	"os"

	"github.com/slack-go/slack"
)

func SlackFraudAlert(channel string, headline string, msg string, approveLink string) {
	fileName := "slack_fraud_slert"
	emptyErrMap := make(map[string]interface{})

	OAUTH_TOKEN := os.Getenv("SLACK_BEARER_TOKEN")
	CHANNEL_ID := channel

	var headerText *slack.TextBlockObject
	headerString := "*" + headline + "*"
	headerText = slack.NewTextBlockObject("mrkdwn", headerString, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	statsField := slack.NewTextBlockObject("mrkdwn", msg, false, false)
	fieldSlice1 := make([]*slack.TextBlockObject, 0)
	fieldSlice1 = append(fieldSlice1, statsField)
	rowSection1 := slack.NewSectionBlock(nil, fieldSlice1, nil)

	api := slack.New(OAUTH_TOKEN)

	previewMessage := headline

	attachment := slack.Attachment{
		Text: "Click yes to ban this domain.",
		Actions: []slack.AttachmentAction{
			{
				Type: "button",
				Text: "Yes",
				URL:  approveLink,
			},
		},
	}

	channelId, timestamp, err := api.PostMessage(
		CHANNEL_ID,
		slack.MsgOptionText(previewMessage, false),
		slack.MsgOptionBlocks(
			headerSection,
			rowSection1,
		),
		slack.MsgOptionAsUser(true),
		slack.MsgOptionAttachments(attachment),
	)

	if err == nil {
		strErr := fmt.Sprintf("Message successfully sent to Channel %s at %s\n", channelId, timestamp)
		logger.LogError("INFO", fileName, nil, strErr, emptyErrMap)
	} else {
		logger.LogError("ERROR", fileName, err, "Error Sending Slack Message", emptyErrMap)
	}
}
