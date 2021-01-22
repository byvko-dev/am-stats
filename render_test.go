package main

import (
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/cufee/am-stats/config"
	dataprep "github.com/cufee/am-stats/dataprep/achievements"
	"github.com/cufee/am-stats/handlers"
	mongodbapi "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	render "github.com/cufee/am-stats/render/achievements"
	"github.com/fogleman/gg"
)

func TestPlayerAchievementsLbImage(t *testing.T) {
	var request handlers.AchievementsRequest
	request.Realm = "NA"
	request.Highlight = false
	request.Days = 0
	request.Limit = 10
	request.PlayerID = 0
	request.Medals = []mongodbapi.MedalWeight{{Name: "markofmastery", Weight: 25}, {Name: "markofmasteryi", Weight: 5}, {Name: "markofmasteryii", Weight: 1}, {Name: "markofmasteryiii", Weight: 1}}

	// Get data
	data, checkData, err := dataprep.ExportAchievementsLeaderboard(request.Realm, request.Days, request.Limit, request.PlayerID, request.Medals...)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}
	if len(data) == 0 {
		log.Print("No data to render")
		t.FailNow()
		return
	}

	// Get BG
	bgImage, err := gg.LoadImage(config.AssetsPath + config.DefaultBG)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}

	// Render image
	image, err := render.PlayerAchievementsLbImage(data, checkData, bgImage, request.Medals)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}

	// Open file
	f, _ := os.Create("test.png")
	defer f.Close()

	// Encode image
	err = png.Encode(f, image)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}
}

func TestClanAchievementsLbImage(t *testing.T) {
	var request handlers.AchievementsRequest
	request.Realm = "NA"
	request.Highlight = false
	request.Days = 0
	request.Limit = 10
	request.ClanTag = "RGN"
	request.Medals = []mongodbapi.MedalWeight{{Name: "markofmastery", Weight: 25}, {Name: "markofmasteryi", Weight: 5}, {Name: "markofmasteryii", Weight: 1}, {Name: "markofmasteryiii", Weight: 1}}

	// Get data
	data, err := dataprep.ExportClanAchievementsLbByRealm(request.Realm, request.Days, request.Limit, request.Medals...)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}
	if len(data) == 0 {
		log.Print("No data to render")
		t.FailNow()
		return
	}

	// Get BG
	bgImage, err := gg.LoadImage(config.AssetsPath + config.DefaultBG)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}

	// Render image
	image, err := render.ClanAchievementsLbImage(data, bgImage, request.Medals)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}

	// Open file
	f, _ := os.Create("test.png")
	defer f.Close()

	// Encode image
	err = png.Encode(f, image)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}
}
