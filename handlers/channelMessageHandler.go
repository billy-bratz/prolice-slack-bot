package handlers

import (
	"fmt"
	"os"
	"prolice-slack-bot/posts"
	"prolice-slack-bot/types"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"mvdan.cc/xurls/v2"
)

func HandleMessageEvent(event *slackevents.MessageEvent, client *slack.Client, currentPRs *[]types.PullRequest) error {

	godotenv.Load(".env")

	matchingString := os.Getenv("MESSAGE_MATCHING_STR")
	secondMatchingString := os.Getenv("MESSAGE_MATCHING_STR2")
	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}

	text := strings.ToLower(event.Text)

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

	if strings.Contains(text, matchingString) &&
		strings.Contains(text, secondMatchingString) {

		xurlsStrict := xurls.Strict()
		prFromText := xurlsStrict.FindAllString(text, -1)
		pr := types.PullRequest{User: fmt.Sprintf("%s %s", user.Profile.FirstName, user.Profile.LastName),
			PrUrl:  prFromText[len(prFromText)-1],
			Posted: time.Now()}

		*currentPRs = append(*currentPRs, pr)

		attachment.Text = fmt.Sprintf("PR: %s added", pr.PrUrl)
		attachment.Color = "#4af030"

		posts.PostMessageWithErrorLogging(client.PostMessage, event.Channel, slack.MsgOptionAttachments(attachment))
	}

	//attachment.Text = "Its Midnight"

	if time.Now().Hour() == 0 {
		posts.PostMessageWithErrorLogging(client.PostMessage, event.Channel, slack.MsgOptionAttachments(attachment))
	}

	return nil
}
