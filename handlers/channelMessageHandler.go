package handlers

import (
	"fmt"
	"log"
	"os"
	"prolice-slack-bot/extensions"
	"prolice-slack-bot/gateways"
	"prolice-slack-bot/posts"
	"prolice-slack-bot/types"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"mvdan.cc/xurls/v2"
)

func HandleMessageEvent(event *slackevents.MessageEvent, client *slack.Client, currentPRs *[]types.PullRequest, hasPosted *bool, silenced *bool) error {

	godotenv.Load(".env")
	matchingString := os.Getenv("MESSAGE_MATCHING_STR")
	secondMatchingString := os.Getenv("MESSAGE_MATCHING_STR2")

	blue := "#2500E0"
	red := "#E31E33"
	currentColor := red

	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		log.Print(err)
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
		strings.Contains(text, secondMatchingString) &&
		!strings.Contains(text, "remove") {

		xurlsStrict := xurls.Strict()
		prMatch := xurlsStrict.FindAllString(text, -1)

		attachment.Text = emptyString

		for _, p := range prMatch {

			if !strings.Contains(p, matchingString) {
				continue
			}

			r, _ := regexp.Compile("(\\d+)")
			prId := r.FindString(p)

			if prId == "" {
				continue
			}

			prExists := extensions.Contains(*currentPRs,
				func(c types.PullRequest) bool {
					return c.Id == prId
				})

			if prExists {
				continue
			}

			pr := gateways.PullRequestById(prId)

			if pr.IsDraft || !strings.EqualFold(pr.Status, "active") {

				log.Printf("Draft or Inactive Pull Request posted, %s\n", p)

				attachment.Text = fmt.Sprintf("Unable to track Pull Request: %s\nPull Requests must be active, and published\n", p)
				attachment.Color = "#EED202"
				PostEphemeral(client, event.Channel, user.ID, *silenced, slack.MsgOptionAttachments(attachment))
				continue
			}

			loc, _ := time.LoadLocation("America/Phoenix")
			pr.PrUrl = p
			pr.Posted = time.Now().In(loc)
			pr.Id = prId

			*currentPRs = append(*currentPRs, pr)

			attachment.Text = fmt.Sprintf("Now tracking Pull Request: %s\n", pr.PrUrl)
			attachment.Color = currentColor

			if !*silenced && len(attachment.Text) > 0 {
				posts.PostMessageWithErrorLogging(client.PostMessage, event.Channel, slack.MsgOptionAttachments(attachment))
			}

			if currentColor == red {
				currentColor = blue
			} else {
				currentColor = red
			}
		}
	}

	if (time.Now().Hour() == 9 || time.Now().Hour() == 12 || time.Now().Hour() == 15) && !*silenced {

		if !*hasPosted {

			for i, pr := range *currentPRs {
				prCheckResult := gateways.PullRequestById(pr.Id)
				log.Printf(prCheckResult.Status)
				if !strings.EqualFold(prCheckResult.Status, "active") {
					*currentPRs = append((*currentPRs)[:i], (*currentPRs)[i+1:]...)
				}
			}
			for _, pr := range *currentPRs {
				attachment.Text = ""
				attachment.Color = currentColor

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
				attachment.Text += fmt.Sprintf("Uncompleted PR by: %s\nUrl: %s\n%s", pr.User, pr.PrUrl, reviewers)
				if attachment.Text != "" {
					posts.PostMessageWithErrorLogging(client.PostMessage, event.Channel, slack.MsgOptionAttachments(attachment))
				}
				if currentColor == red {
					currentColor = blue
				} else {
					currentColor = red
				}
			}

			*hasPosted = true
		}
	} else {
		*hasPosted = false
	}

	return nil
}

func PostEphemeral(client *slack.Client, channelId string, userId string, silenced bool, options slack.MsgOption) {
	if !silenced {
		client.PostEphemeral(channelId, userId, options)
	}
}
