package handlers

import (
	"errors"
	"prolice-slack-bot/types"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func HandleChannelMessage(event slackevents.EventsAPIEvent, client *slack.Client, currentPRs *[]types.PullRequest, hasPosted *bool) error {
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
			err := HandleMessageEvent(ev, client, currentPRs, hasPosted)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}
