package handlers

import (
	"log"
	"net/http"
	"runtime/debug"

	achievements "github.com/cufee/am-stats/dataprep/achievements"
	"github.com/gofiber/fiber/v2"
)

// HandleAchievementsJSONExport - Get achievements as JSON
func HandleAchievementsJSONExport(c *fiber.Ctx) error {
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