package render

import (
	"testing"

	"github.com/cufee/am-stats/stats"
)

func TestChallengeRender(t *testing.T) {
	pid := 1013072123
	// pid := 1023629188
	realm := "NA"
	export, _ := stats.ExportSessionAsStruct(pid, realm, 0)
	_ = export
}
