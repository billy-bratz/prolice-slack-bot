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

func HandleMessageEvent(event *slackevents.MessageEvent, client *slack.Client, currentPRs *[]types.PullRequest, hasPosted *bool) error {

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
		strings.Contains(text, secondMatchingString) &&
		!strings.Contains(text, "remove") {

		xurlsStrict := xurls.Strict()
		prMatch := xurlsStrict.FindAllString(text, -1)
		prUrl := prMatch[len(prMatch)-1]

		r, _ := regexp.Compile("(\\d+)")
		prId := r.FindString(prUrl)

		if prId == "" {
			log.Printf("Could not parse id from url: %s", prUrl)

			attachment.Text = fmt.Sprintf("Could not parse id from url: %s.\nTo track please ensure the pull request id is visible in plain text", prUrl)
			attachment.Color = "#4af030"
			client.PostEphemeral(event.Channel, user.ID, slack.MsgOptionAttachments(attachment))

			return nil
		}

		prExists := extensions.Contains(*currentPRs,
			func(c types.PullRequest) bool {
				return c.PrUrl == prUrl
			})

		if prExists {
			log.Printf("Pull Request already exists: %s", prUrl)

			attachment.Text = fmt.Sprintf("Already tracking Pull Request: %s", prUrl)
			attachment.Color = "#4af030"
			client.PostEphemeral(event.Channel, user.ID, slack.MsgOptionAttachments(attachment))

			return nil
		}

		pr := gateways.PullRequestById(prId)

		if pr.IsDraft || !strings.EqualFold(pr.Status, "active") {

			log.Printf("Draft or Inactive Pull Request posted, %s", prUrl)

			attachment.Text = fmt.Sprintf("Unable to track Pull Request: %s\nPull Requests must be active, and published", prUrl)
			attachment.Color = "#4af030"
			client.PostEphemeral(event.Channel, user.ID, slack.MsgOptionAttachments(attachment))

			return nil
		}

		loc, _ := time.LoadLocation("America/Phoenix")
		pr.PrUrl = prUrl
		pr.Posted = time.Now().In(loc)
		pr.Id = prId

		*currentPRs = append(*currentPRs, pr)

		attachment.Text = fmt.Sprintf("Now tracking Pull Request: %s", pr.PrUrl)
		attachment.Color = "#4af030"

		posts.PostMessageWithErrorLogging(client.PostMessage, event.Channel, slack.MsgOptionAttachments(attachment))
	}

	if time.Now().Hour() == 9 || time.Now().Hour() == 12 || time.Now().Hour() == 15 {

		if !*hasPosted {

			for i, pr := range *currentPRs {
				attachment.Text = ""
				prCheckResult := gateways.PullRequestById(pr.Id)
				log.Printf(prCheckResult.Status)
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

					attachment.Text += fmt.Sprintf("Uncompleted PR by: %s\nUrl: %s\n%s", pr.User, pr.PrUrl, reviewers)

					if attachment.Text != "" {
						posts.PostMessageWithErrorLogging(client.PostMessage, event.Channel, slack.MsgOptionAttachments(attachment))
					}
				}
			}

			*hasPosted = true
		}
	} else {
		*hasPosted = false
	}

	return nil
}
