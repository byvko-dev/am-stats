package dataprep

import (
	"log"
	"testing"

	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
)

func TestExportClanAchievementsByTag(t *testing.T) {
	medals := []dbAch.MedalWeight{{Name: "markofmastery", Weight: 1}}
	clanTag := "RGN"
	realm := "NA"

	data, clanTotal, err := ExportClanAchievementsByTag(clanTag, realm, 0, medals...)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}
	log.Printf("%+v", data)
	log.Printf("%v", clanTotal)
}
