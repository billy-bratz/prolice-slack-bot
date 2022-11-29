package contract

import "prolice-slack-bot/types"

type PullRequestResponse struct {
	Status    string
	IsDraft   bool
	Reviewers []Reviewer
	CreatedBy CreatedBy
}

type Reviewer struct {
	Vote        int
	DisplayName string
	UniqueName  string
}

type CreatedBy struct {
	DisplayName string
	UniqueName  string
}

func ToType(prResponse PullRequestResponse) types.PullRequest {
	return types.PullRequest{
		Status:    prResponse.Status,
		User:      prResponse.CreatedBy.DisplayName,
		Reviewers: AddReviewers(prResponse.Reviewers),
		IsDraft:   prResponse.IsDraft,
	}
}

func AddReviewers(reviewers []Reviewer) []types.Reviewer {
	var returnReviewers []types.Reviewer

	for _, r := range reviewers {
		rvr := types.Reviewer{
			Vote:        r.Vote,
			DisplayName: r.DisplayName,
			UniqueName:  r.UniqueName,
		}
		returnReviewers = append(returnReviewers, rvr)
	}

	return returnReviewers
}
