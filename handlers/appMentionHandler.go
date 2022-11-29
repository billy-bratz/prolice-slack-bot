package handlers

import (
	"fmt"
	"prolice-slack-bot/extensions"
	"prolice-slack-bot/gateways"
	"prolice-slack-bot/types"
	"strings"

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
	if strings.Contains(text, "list prs") {

		for i, pr := range *currentPRs {
			prCheckResult := gateways.PullRequestById(pr.Id)

			if !strings.EqualFold(prCheckResult.Status, "active") {
				*currentPRs = append((*currentPRs)[:i], (*currentPRs)[i+1:]...)
			} else {
				reviewers := "Approvals: "
				for _, r := range pr.Reviewers {
					if r.Vote == 5 && !strings.Contains(r.DisplayName, "[CarvanaDev]") {
						reviewers += fmt.Sprintf("%s approved with suggestions\n", r.DisplayName)
					} else if r.Vote == 10 && !strings.Contains(r.DisplayName, "[CarvanaDev]") {
						reviewers += fmt.Sprintf("%s approved\n", r.DisplayName)
					}
				}

				if len(reviewers) < 13 {
					reviewers += "None"
				}
				attachment.Text = fmt.Sprintf("\nUrl: %s", pr.PrUrl)
				attachment.Text += fmt.Sprintf("\nAuthor: %s", pr.User)
				attachment.Text += fmt.Sprintf("\nPosted Date: %s", pr.Posted)
				attachment.Text += fmt.Sprintf("\n%s", reviewers)

				if attachment.Text == "" {
					attachment.Text = "No PRs"
				}
			}
		}
	} else if strings.Contains(text, "remove pr") {
		xurlsStrict := xurls.Strict()
		prMatch := xurlsStrict.FindAllString(text, -1)
		prUrl := prMatch[len(prMatch)-1]

		i := extensions.IndexOf(*currentPRs, func(pr types.PullRequest) bool { return pr.PrUrl == prUrl })
		*currentPRs = append((*currentPRs)[:i], (*currentPRs)[i+1:]...)
		attachment.Text = ("PR removed")
	} else {
		attachment.Text = fmt.Sprintf("Sorry %s, I do not know how to handle that request", user.Name)
	}

	attachment.Color = "#4af030"
	client.PostEphemeral(event.Channel, user.ID, slack.MsgOptionAttachments(attachment))

	return nil
}
