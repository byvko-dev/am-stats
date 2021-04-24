package main

import (
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/cufee/am-stats/config"
	dataprep "github.com/cufee/am-stats/dataprep/achievements"
	replays "github.com/cufee/am-stats/dataprep/replays"
	stats "github.com/cufee/am-stats/dataprep/stats"
	"github.com/cufee/am-stats/handlers"
	mongodbapi "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	render "github.com/cufee/am-stats/render/achievements"
	renderReplay "github.com/cufee/am-stats/render/replays"
	renderStats "github.com/cufee/am-stats/render/stats"
	"github.com/fogleman/gg"
)

func TestPlayerAchievementsLbImage(t *testing.T) {
	var request handlers.AchievementsRequest
	request.Days = 0
	request.Limit = 10
	request.Realm = "NA"
	request.PlayerID = 0
	request.Highlight = false
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
	data, check, err := dataprep.ExportClanAchievementsLbByRealm(request.Realm, request.PlayerID, request.Days, request.Limit, request.Medals...)
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
	image, err := render.ClansAchievementsLbImage(data, check, bgImage, request.Medals)
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

func TestReplayRender(t *testing.T) {
	url := "https://cdn.discordapp.com/attachments/346861614735294464/821871066342752256/20210317_1816__Mochanado_Cz04_T50_51_4295187955351455.wotbreplay"
	// url := "https://cdn.discordapp.com/attachments/346861614735294464/816882621723443221/20210303_1914__Tutankhamun_1332BC_A116_XM551_11598169955651744.wotbreplay"
	// url := "https://cdn.discordapp.com/attachments/719831153162321981/823314368333873192/20210320_2053__Vladok1408_Object252_4311306967881970.wotbreplay"

	// Get data
	data, err := replays.ProcessReplay(url)
	if err != nil {
		log.Print(err)
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
	image, err := renderReplay.Render(data, bgImage)
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
func TestStatsRender(t *testing.T) {
	realm := "RU"
	id := 112833260
	days := 30

	// Get data
	data, err := stats.ExportSessionAsStruct(id, 0, realm, days, 3, "", false, true)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}

	// Get BG
	bgImage, err := gg.LoadImage(config.AssetsPath + "bg_code_fatal.jpg")
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}

	// Render image
	image, err := renderStats.ImageFromStats(data, "", 3, false, true, bgImage)
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
