package helpers

import (
	"prolice-slack-bot/types"
	"strings"
)

func RemoveInactivePrs(status string, idx int, currentPrs *[]types.PullRequest) bool {

	if !strings.EqualFold(status, "active") {
		RemovePr(idx, currentPrs)
		return true
	}

	return false
}

func RemovePr(idx int, currentPrs *[]types.PullRequest) {
	*currentPrs = append((*currentPrs)[:idx], (*currentPrs)[idx+1:]...)
}
