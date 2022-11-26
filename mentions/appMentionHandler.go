package mentions

import (
	"fmt"
	"prolice-slack-bot/posts"
	"prolice-slack-bot/types"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"mvdan.cc/xurls/v2"
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
	if strings.Contains(text, "https://carvanadev.visualstudio.com/carvana.underwriting/") && strings.Contains(text, "pullrequest") {

		xurlsStrict := xurls.Strict()
		prFromText := xurlsStrict.FindAllString(text, -1)
		pr := types.PullRequest{User: fmt.Sprintf("%s %s", user.Profile.FirstName, user.Profile.LastName), PrUrl: prFromText[len(prFromText)-1], Posted: time.Now()}

		*currentPRs = append(*currentPRs, pr)

		attachment.Text = fmt.Sprintf("PR: %s added", pr.PrUrl)
		attachment.Color = "#4af030"

		posts.PostMessageWithErrorLogging(client.PostMessage, event.Channel, slack.MsgOptionAttachments(attachment))

		client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	} else if strings.Contains(text, "list prs") {
		// Send a message to the user

		for _, pr := range *currentPRs {
			//use attachment.Fields here
			attachment.Text = fmt.Sprintf("\nUrl: %s", pr.PrUrl)
			attachment.Text += fmt.Sprintf("\nAuthor: %s", pr.User)
			attachment.Text += fmt.Sprintf("\nPosted Date: %s", pr.Posted)

			client.PostEphemeral(event.Channel, user.ID, slack.MsgOptionAttachments(attachment))
		}

		// attachment.Pretext = "How can I be of service"
		attachment.Color = "#4af030"
	} else if strings.Contains(text, "trigger") {

		attachment.Text = "triggered"
		posts.PostMessageWithErrorLogging(client.PostMessage, event.Channel, slack.MsgOptionAttachments(attachment))
	} else {
		// Send a message to the user
		attachment.Text = fmt.Sprintf("I am good. How are you %s?", user.Name)
		// attachment.Pretext = "How can I be of service"
		attachment.Color = "#4af030"
	}

	return nil
}
