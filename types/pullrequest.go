package types

import "time"

type PullRequest struct {
	User      string
	PrUrl     string
	Posted    time.Time
	Reviewers []Reviewer
	Status    string
	IsDraft   bool
	Id        string
}

type Reviewer struct {
	Vote        int
	DisplayName string
	UniqueName  string
}
