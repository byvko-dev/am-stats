package handlers

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/cufee/am-stats/config"
	stats "github.com/cufee/am-stats/dataprep/stats"
	statsRender "github.com/cufee/am-stats/render/stats"
	externalapis "github.com/cufee/am-stats/wargamingapi"
	"github.com/fogleman/gg"
	"github.com/gofiber/fiber/v2"
)

// HandleStatsImageExport - Export stats as image
func HandleStatsImageExport(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	var request StatsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Log player ID and realm
	log.Printf("pid: %v, realm: %v", request.PlayerID, request.Realm)

	export, err := stats.ExportSessionAsStruct(request.PlayerID, request.TankID, request.Realm, request.Days, 0, "", request.Special, request.IncludeRating)
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
		bgImage, err = gg.LoadImage(config.AssetsPath + config.DefaultBG)
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("failed to load a background image: %#v", err),
			})
		}
	}

	img, err := statsRender.ImageFromStats(export, request.Sort, request.TankLimit, request.Premium, request.Verified, bgImage, export.PlayerCache.UniquePins...)
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

// HandleStatsJSONExport - Get stats as JSON
func HandleStatsJSONExport(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	var request StatsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	export, err := stats.ExportSessionAsStruct(request.PlayerID, request.TankID, request.Realm, request.Days, request.TankLimit, request.Sort, request.Special, request.IncludeRating)
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

// HandleSpecialSessionReset - reset special session
func HandleSpecialSessionReset(c *fiber.Ctx) error {
	var request StatsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	httpClient := &http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	// Make request
	req, err := http.NewRequest("GET", config.AMCacheURL+fmt.Sprintf("/%v/special-session/%v", request.Realm, request.PlayerID), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.AMAPIKey)

	// Send request
	res, err := httpClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check response code
	if res.StatusCode != fiber.StatusOK {
		errData, _ := ioutil.ReadAll(res.Body)
		return c.Status(res.StatusCode).Send(errData)
	}
	return c.SendStatus(fiber.StatusOK)
}

// HandlePublicStatsJSONExport - Get stats as JSON
func HandlePublicStatsJSONExport(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	var request StatsRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	export, err := stats.ExportSessionAsStruct(request.PlayerID, request.TankID, request.Realm, request.Days, request.TankLimit, request.Sort, request.Special, request.IncludeRating)

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

	// Trim data
	var publicExport stats.ExportData
	publicExport.SessionStats = export.SessionStats
	publicExport.LastSession = export.LastSession
	publicExport.PlayerDetails = export.PlayerDetails

	return c.JSON(publicExport)
}
