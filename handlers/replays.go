package handlers

import (
	"log"
	"net/http"
	"runtime/debug"

	replays "github.com/cufee/am-stats/dataprep/replays"
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

	// Get replay URL
	replayURL := c.Query("url")
	if replayURL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid replay url",
		})
	}

	// Export data
	export, err := replays.ProcessReplay(replayURL)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(export)
}
