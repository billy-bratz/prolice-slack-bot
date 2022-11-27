package handlers

import (
	"fmt"
	"prolice-slack-bot/types"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleAppMentionEventToBot is used to take care of the AppMentionEvent when the bot is mentioned
func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client, currentPRs *[]types.PullRequest) error {

	//TODO//
	// Check the PR's we have and if they're still valid
	// remove the ones that are not, schedule a message for the next one

	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	// Check if the user said Hello to the bot
	text := strings.ToLower(event.Text)

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
	if strings.Contains(text, "list prs") {

		for _, pr := range *currentPRs {
			//use attachment.Fields here
			attachment.Text = fmt.Sprintf("\nUrl: %s", pr.PrUrl)
			attachment.Text += fmt.Sprintf("\nAuthor: %s", pr.User)
			attachment.Text += fmt.Sprintf("\nPosted Date: %s", pr.Posted)
		}
	} else {
		attachment.Text = fmt.Sprintf("Sorry %s, I do not know how to handle that request", user.Name)
		// attachment.Pretext = "How can I be of service"
	}

	attachment.Color = "#4af030"
	client.PostEphemeral(event.Channel, user.ID, slack.MsgOptionAttachments(attachment))

	return nil
}
