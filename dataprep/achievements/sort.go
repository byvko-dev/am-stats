package dataprep

import (
	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	"github.com/jfcg/sorty"
)

func SortPlayerLeaderboard(board []dbAch.AchievementsPlayerData) {
	lsw := func(i, k, r, s int) bool {
		if board[i].Score > board[k].Score {
			if r != s {
				board[r], board[s] = board[s], board[r]
			}
			return true
		}
		return false
	}
	sorty.Sort(len(board), lsw)
}
