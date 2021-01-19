package handlers

import (
	"net/http"
	"strconv"

	"github.com/cufee/am-stats/mongodbapi"
	externalapis "github.com/cufee/am-stats/wargamingapi"
	"github.com/gofiber/fiber/v2"
)

// HandlePlayerCheck - Get player info by id
func HandlePlayerCheck(c *fiber.Ctx) error {
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
