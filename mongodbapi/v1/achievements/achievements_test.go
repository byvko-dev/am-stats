package mongodbapi

import (
	"log"
	"testing"
)

func TestGetPlayerAchievementsLb(t *testing.T) {
	fields := []string{"data.achievements.markofmastery", "data.achievements.markofmasteryi", "data.achievements.markofmasteryii", "data.achievements.markofmasteryiii"}
	data, err := GetPlayerAchievementsLb("EU", fields...)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}
	log.Printf("%+v", data)
}
