package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/cufee/am-stats/config"
	achievements "github.com/cufee/am-stats/dataprep/achievements"
	dataprep "github.com/cufee/am-stats/dataprep/achievements"
	render "github.com/cufee/am-stats/render/achievements"
	"github.com/fogleman/gg"
	"github.com/gofiber/fiber/v2"
)

// HandleSessionAchievementsExport - Get achievements as JSON
func HandleSessionAchievementsExport(c *fiber.Ctx) error {
	// Recover on panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	// Parse request data
	var request StatsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get data
	export, err := achievements.ExportAchievementsSession(request.PlayerID, request.Realm, request.Days)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(export)
}

// HandleClanAchievementsExport - Get achievements Leaderboard as JSON
func HandleClanAchievementsExport(c *fiber.Ctx) error {
	// Recover on panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	// Parse request data
	var request AchievementsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check request
	if request.ClanTag == "" && request.Realm == "" {
		return fiber.ErrBadRequest
	}
	if len(request.Medals) < 1 {
		return fiber.ErrBadRequest
	}

	log.Print(request.ClanTag, request.Realm, request.Medals)

	// Get data
	export, total, err := achievements.ExportClanAchievementsByTag(request.ClanTag, request.Realm, request.Days, request.Medals...)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"leaderboard": export,
		"total_score": total,
	})
}

// HandleClanAchievementsLbExport - Get achievements Leaderboard as JSON
func HandleClanAchievementsLbExport(c *fiber.Ctx) error {
	// Recover on panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	// Parse request data
	var request AchievementsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check request
	if request.ClanTag == "" && request.Realm == "" {
		return fiber.ErrBadRequest
	}
	if len(request.Medals) < 1 {
		return fiber.ErrBadRequest
	}

	// Get data
	export, checkData, err := achievements.ExportClanAchievementsLbByRealm(request.Realm, request.PlayerID, request.Days, request.Limit, request.Medals...)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"leaderboard": export,
		"player_clan": checkData,
	})
}

// HandlePlayersAchievementsLbExport - Get achievements Leaderboard as JSON
func HandlePlayersAchievementsLbExport(c *fiber.Ctx) error {
	// Recover on panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	// Parse request data
	var request AchievementsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get data
	export, position, err := achievements.ExportAchievementsLeaderboard(request.Realm, request.Days, request.Limit, request.PlayerID, request.Medals...)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"leaderboard":     export,
		"player_position": position,
	})
}

// HandlerPlayersLeaderboardImage -
func HandlerPlayersLeaderboardImage(c *fiber.Ctx) error {
	// Recover on panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	// Parse request data
	var request AchievementsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get bg Image
	var bgImage image.Image
	if request.BgURL != "" {
		response, _ := http.Get(request.BgURL)
		if response != nil {
			bgImage, _, err = image.Decode(response.Body)
			defer response.Body.Close()
		} else {
			log.Printf("bad bg image for %v", request.PlayerID)
			err = fmt.Errorf("bad bg image")
		}
	}
	if err != nil || request.BgURL == "" {
		bgImage, err = gg.LoadImage(config.AssetsPath + config.DefaultBG)
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("failed to load a background image: %#v", err),
			})
		}
	}

	// Get data
	data, checkData, err := dataprep.ExportAchievementsLeaderboard(request.Realm, request.Days, request.Limit, request.PlayerID, request.Medals...)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if checkData.Position > request.Limit {
		data = append(data, checkData.AchievementsPlayerData)
	}

	// Render image
	image, err := render.PlayerAchievementsLbImage(data, checkData, bgImage, request.Medals)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Encode image
	buf := new(bytes.Buffer)
	err = png.Encode(buf, image)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Send image
	c.Set("Content-Type", "image/png")
	s, err := ioutil.ReadAll(buf)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Send(s)
}

// HandlerClansLeaderboardImage -
func HandlerClansLeaderboardImage(c *fiber.Ctx) error {
	// Recover on panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	// Parse request data
	var request AchievementsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get bg Image
	var bgImage image.Image
	if request.BgURL != "" {
		response, _ := http.Get(request.BgURL)
		if response != nil {
			bgImage, _, err = image.Decode(response.Body)
			defer response.Body.Close()
		} else {
			log.Printf("bad bg image for %v", request.PlayerID)
			err = fmt.Errorf("bad bg image")
		}
	}
	if err != nil || request.BgURL == "" {
		bgImage, err = gg.LoadImage(config.AssetsPath + config.DefaultBG)
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("failed to load a background image: %#v", err),
			})
		}
	}

	// Check request
	if request.ClanTag == "" && request.Realm == "" {
		return fiber.ErrBadRequest
	}
	if len(request.Medals) < 1 {
		return fiber.ErrBadRequest
	}

	// Get data
	data, checkData, err := achievements.ExportClanAchievementsLbByRealm(request.Realm, request.PlayerID, request.Days, request.Limit, request.Medals...)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Add player clan to data
	if checkData.Position > request.Limit {
		data = append(data, checkData)
	}

	// Render image
	image, err := render.ClansAchievementsLbImage(data, checkData, bgImage, request.Medals)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Encode image
	buf := new(bytes.Buffer)
	err = png.Encode(buf, image)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Send image
	c.Set("Content-Type", "image/png")
	s, err := ioutil.ReadAll(buf)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Send(s)
}
