package handlers

import (
	"fmt"
	"log"
	"prolice-slack-bot/extensions"
	"prolice-slack-bot/gateways"
	"prolice-slack-bot/helpers"
	"prolice-slack-bot/types"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"mvdan.cc/xurls/v2"
)

const listPrs = "list"
const removePr = "remove"
const emptyString = ""
const silence = "hush"
const unmute = "unmute"
const MMDDYYY = "01-02-2006"

// HandleAppMentionEventToBot is used to take care of the AppMentionEvent when the bot is mentioned
func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client, currentPRs *[]types.PullRequest, silenced *bool) error {

	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		log.Print(err)
	}

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
	switch {
	case strings.Contains(text, listPrs):
		for i, pr := range *currentPRs {
			prCheckResult := gateways.PullRequestById(pr.Id)

			wasRemoved := helpers.RemoveInactivePrs(prCheckResult.Status, i, *&currentPRs)

			if !wasRemoved {
				attachment.Text += BuildPrList(pr)
			}
		}

		if len(attachment.Text) == 0 {
			attachment.Text = "No Pull Requests for review"
		}

	case strings.Contains(text, removePr):
		xurlsStrict := xurls.Strict()
		prMatch := xurlsStrict.FindAllString(text, -1)

		if len(prMatch) == 0 {
			attachment.Text = "Could not find URL in string please use the syntax remove (url)"
		} else {
			prUrl := prMatch[len(prMatch)-1]
			i := extensions.IndexOf(*currentPRs, func(pr types.PullRequest) bool { return pr.PrUrl == prUrl })
			if i == -1 {
				attachment.Text = "Could not find matching Pull Request to remove.\n"
				attachment.Text += "Use `@PRolice list` to see list of Pull Requests."
			} else {
				helpers.RemovePr(i, *&currentPRs)
				attachment.Text = "Pull Request removed"
			}
		}
	case strings.Contains(text, silence):

		if *silenced {
			attachment.Text = "already silenced do you wish to use unmute?"
			PostResponse(attachment, client, event, user)
		} else {
			attachment.Text = "going dark"
			PostResponse(attachment, client, event, user)
			*silenced = true
		}
	case strings.Contains(text, unmute):
		if *silenced {
			attachment.Text = "unmuted"
		}
		*silenced = false
	default:
		attachment.Text = fmt.Sprintf("Available Commands:\n")
		attachment.Fields = []slack.AttachmentField{
			{
				Title: "List all Pull Requests",
				Value: fmt.Sprintf("@PRolice %s", listPrs),
			},
			{
				Title: "Remove Pull Request",
				Value: fmt.Sprintf("@PRolice %s", removePr),
			},
			{
				Title: "Silence Chat notifications",
				Value: fmt.Sprintf("@PRolice %s", silence),
			},
		}
	}

	if !*silenced {
		PostResponse(attachment, client, event, user)
	}

	return nil
}

func PostResponse(attachment slack.Attachment, client *slack.Client, event *slackevents.AppMentionEvent, user *slack.User) {
	attachment.Color = "#4af030"
	client.PostEphemeral(event.Channel, user.ID, slack.MsgOptionAttachments(attachment))
}

func BuildPrList(pullRequest types.PullRequest) string {
	response := emptyString
	reviewers := "Approvals: "

	for _, r := range pullRequest.Reviewers {
		if r.Vote == 5 && !strings.Contains(r.DisplayName, "[CarvanaDev]") {
			reviewers += fmt.Sprintf("%s approved with suggestions\n", r.DisplayName)
		} else if r.Vote == 10 && !strings.Contains(r.DisplayName, "[CarvanaDev]") {
			reviewers += fmt.Sprintf("%s approved\n", r.DisplayName)
		}
	}

	if len(reviewers) < 13 {
		reviewers += "None"
	}

	response = fmt.Sprintf("\nUrl: %s", pullRequest.PrUrl)
	response += fmt.Sprintf("\nAuthor: %s", pullRequest.User)
	response += fmt.Sprintf("\nPosted Date: %s", pullRequest.Posted.Format(MMDDYYY))
	response += fmt.Sprintf("\n%s ", reviewers)

	return response
}
