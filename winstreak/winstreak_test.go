package winstreak

import (
	"log"
	"testing"

	"github.com/cufee/am-stats/stats"
)

func TestWinStreak(t *testing.T) {
	pid := 1013072123
	realm := "NA"
	export, _ := stats.ExportSessionAsStruct(pid, realm, 0)
	streak, err := CheckStreak(export.PlayerDetails.ID, export.PlayerDetails.Stats.All)
	log.Print(streak, err)
}
