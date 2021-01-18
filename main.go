package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"runtime/debug"
	"strconv"

	"github.com/cufee/am-stats/auth"
	"github.com/cufee/am-stats/config"
	"github.com/cufee/am-stats/mongodbapi"
	"github.com/cufee/am-stats/render"
	"github.com/cufee/am-stats/stats"
	externalapis "github.com/cufee/am-stats/wargamingapi"
	"github.com/fogleman/gg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"net/http"
)

type request struct {
	PlayerID  int    `json:"player_id"`
	Premium   bool   `json:"premium"`
	Verified  bool   `json:"verified"`
	Realm     string `json:"realm"`
	Days      int    `json:"days"`
	Sort      string `json:"sort_key"`
	TankLimit int    `json:"detailed_limit"`
	BgURL     string `json:"bg_url"`
}

const currentBG string = "bg_frame.png"

func main() {
	// Define routes
	app := fiber.New()

	// Logger
	app.Use(logger.New())

	// API key validator
	app.Use(auth.Validator)

	// Checks
	app.Get("/player/id/:id", handlePlayerCheck)

	// Stats
	app.Get("/player", handlePlayerRequest)
	app.Get("/stats", handleStatsRequest)

	log.Print(app.Listen(fmt.Sprintf(":%v", config.APIport)))
}

func handlePlayerCheck(c *fiber.Ctx) error {
	// Get player ID as int
	playerID := c.Params("id")
	pid, err := strconv.Atoi(playerID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get realm
	realm, err := mongodbapi.GetRealmByPID(pid)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get live player name
	playerData, err := externalapis.PlayerProfileData(pid, realm)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return
	return c.JSON(fiber.Map{
		"nickname": playerData.Name,
		"realm":    realm,
	})
}

func handlePlayerRequest(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	var request request
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Log player ID and realm
	log.Printf("pid: %v, realm: %v", request.PlayerID, request.Realm)

	export, err := stats.ExportSessionAsStruct(request.PlayerID, request.Realm, request.Days)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if export.PlayerDetails == (externalapis.PlayerProfile{}) || export.PlayerDetails.Name == "" {
		log.Printf("%+v", request)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "bad player data",
		})
	}
	if request.TankLimit == 0 {
		request.TankLimit = 10
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
		bgImage, err = gg.LoadImage(config.AssetsPath + currentBG)
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("failed to load a background image: %#v", err),
			})
		}
	}

	img, err := render.ImageFromStats(export, request.Sort, request.TankLimit, request.Premium, request.Verified, bgImage)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Encode image
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
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

func handleStatsRequest(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	var request request
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	export, err := stats.ExportSessionAsStruct(request.PlayerID, request.Realm, request.Days)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if export.PlayerDetails == (externalapis.PlayerProfile{}) || export.PlayerDetails.Name == "" {
		log.Printf("%+v", request)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "bad player data",
		})
	}

	return c.JSON(export)
}
