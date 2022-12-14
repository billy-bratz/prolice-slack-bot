package contract

import "prolice-slack-bot/types"

type PullRequestResponse struct {
	Status     string
	IsDraft    bool
	Reviewers  []Reviewer
	CreatedBy  CreatedBy
	Repository Repository
}

type Reviewer struct {
	Vote        int
	DisplayName string
	UniqueName  string
	IsRequired  bool
}

type CreatedBy struct {
	DisplayName string
	UniqueName  string
}

type Repository struct {
	Id   string
	Name string
	Url  string
}

func ToType(prResponse PullRequestResponse) types.PullRequest {
	return types.PullRequest{
		Status:     prResponse.Status,
		User:       prResponse.CreatedBy.DisplayName,
		Reviewers:  AddReviewers(prResponse.Reviewers),
		IsDraft:    prResponse.IsDraft,
		Repository: types.Repository{Id: prResponse.Repository.Id, Name: prResponse.Repository.Name, Url: prResponse.Repository.Url},
	}
}

func AddReviewers(reviewers []Reviewer) []types.Reviewer {
	var returnReviewers []types.Reviewer

	for _, r := range reviewers {
		rvr := types.Reviewer{
			Vote:        r.Vote,
			DisplayName: r.DisplayName,
			UniqueName:  r.UniqueName,
			IsRequired:  r.IsRequired,
		}
		returnReviewers = append(returnReviewers, rvr)
	}

	return returnReviewers
}
