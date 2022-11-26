package posts

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

func PostMessageWithErrorLogging(f func(string, ...slack.MsgOption) (string, string, error), channelID string, options ...slack.MsgOption) {
	_, _, err := f(channelID, options...)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to post message: %w", err))
	}
}
