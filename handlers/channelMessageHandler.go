package handlers

import (
	"prolice-slack-bot/posts"
	"prolice-slack-bot/types"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func HandleMessageEvent(event *slackevents.MessageEvent, client *slack.Client, currentPRs *[]types.PullRequest) error {

	// Grab the user name based on the ID of the one who mentioned the bot
	_, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	// Check if the user said Hello to the bot
	//text := strings.ToLower(event.Text)

	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}
	// Add Some default context like user who mentioned the bot
	// attachment.Fields = []slack.AttachmentField{
	// 	{
	// 		Title: "Date",
	// 		Value: time.Now().String(),
	// 	}, {
	// 		Title: "Initializer",
	// 		Value: user.Name,
	// 	},
	// }

	attachment.Text = "Its Midnight"

	// becareful for infinate loops here (bot triggering itself)
	if time.Now().Hour() == 0 {
		posts.PostMessageWithErrorLogging(client.PostMessage, event.Channel, slack.MsgOptionAttachments(attachment))
	}

	return nil
}
