package types

import "time"

type PullRequest struct {
	User   string
	PrUrl  string
	Posted time.Time
}
