package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"prolice-slack-bot/posts"
	"prolice-slack-bot/types"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"mvdan.cc/xurls/v2"
)

func main() {

	// Load Env variables from .dot file
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	// Create a new client to slack by giving token
	// Set debug to true while developing
	// Also add a ApplicationToken option to the client
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	// go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
	socket := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	var currentPRs []types.PullRequest

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socket *socketmode.Client) {
		// Create a for loop that selects either the context cancellation or the events incomming
		for {
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socket.Events:
				// We have a new Events, let's type switch the event
				// Add more use cases here if you want to listen to other events.
				switch event.Type {
				// handle EventAPI events
				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socket.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch

					//log.Println(eventsAPI) // commenting for event hanndling

					//------------------------------------
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					err := HandleEventMessage(eventsAPI, client, &currentPRs)
					if err != nil {
						// Replace with actual err handeling
						log.Fatal(err)
					}
				}
			}
		}
	}(ctx, client, socket)

	socket.Run()

	// attachment := slack.Attachment{
	// 	Pretext: "Bot Message",
	// 	Text:    "Text",
	// 	Color:   "4af030",
	// 	Fields: []slack.AttachmentField{
	// 		{
	// 			Title: "Date",
	// 			Value: time.Now().String(),
	// 		},
	// 	},
	// }

	// _, timestamp, err := client.PostMessage(
	// 	channelID,
	// 	slack.MsgOptionAttachments(attachment),
	// )

	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Mesage sent at %s", timestamp)

}

func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client, currentPRs *[]types.PullRequest) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := HandleAppMentionEventToBot(ev, client, currentPRs)
			if err != nil {
				return err
			}
		case *slackevents.MessageEvent:
			err := HandleMessageEvent(ev, client, currentPRs)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

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

	// Send the message to the channel
	// The Channel is available in the event message
	// _, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	// if err != nil {
	// 	return fmt.Errorf("failed to post message: %w", err)
	// }
	return nil
}

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
