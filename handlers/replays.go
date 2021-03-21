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
	replays "github.com/cufee/am-stats/dataprep/replays"
	"github.com/cufee/am-stats/external"
	renderReplay "github.com/cufee/am-stats/render/replays"
	"github.com/fogleman/gg"
	"github.com/gofiber/fiber/v2"
)

// HandleReplayJSONExport - Get replay data as JSON
func HandleReplayJSONExport(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	var request ReplayRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if request.ReplayURL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid replay url",
		})
	}

	// Export data
	export, err := replays.ProcessReplay(request.ReplayURL)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(export)
}

// HandleReplayImageExport - Get replay data as Image
func HandleReplayImageExport(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "something did not work",
			})
		}
	}()

	var request ReplayRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if request.ReplayURL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid replay url",
		})
	}

	// Export data
	export, err := replays.ProcessReplay(request.ReplayURL)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get protogonist bg URL
	protoData, _ := external.CheckUserByPID(export.Protagonist)

	// Get bg Image
	var bgImage image.Image
	if protoData.CustomBgURL != "" {
		response, _ := http.Get(protoData.CustomBgURL)
		if response != nil {
			bgImage, _, err = image.Decode(response.Body)
			defer response.Body.Close()
		} else {
			err = fmt.Errorf("bad bg image")
		}
	}
	if err != nil || protoData.CustomBgURL == "" {
		bgImage, err = gg.LoadImage(config.AssetsPath + config.DefaultBG)
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("failed to load a background image: %#v", err),
			})
		}
	}

	// Render image
	image, err := renderReplay.Render(export, bgImage)
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
