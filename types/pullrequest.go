package types

import "time"

type PullRequest struct {
	User       string
	PrUrl      string
	Posted     time.Time
	Reviewers  []Reviewer
	Status     string
	IsDraft    bool
	Id         string
	Repository Repository
}

type Reviewer struct {
	Vote        int
	DisplayName string
	UniqueName  string
	IsRequired  bool
}

type Repository struct {
	Id   string
	Name string
	Url  string
}
