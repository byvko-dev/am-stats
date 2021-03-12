package dataprep

import (
	"log"
	"testing"
)

func TestExportSessionAsStruct(t *testing.T) {
	pid := 1013072123
	// pid := 1023629188
	realm := "NA"
	export, err := ExportSessionAsStruct(pid, 0, realm, 0, 0, "", true)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}
	log.Printf("%+v", export.SessionStats.BattlesAll)
}
