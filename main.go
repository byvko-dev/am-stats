package main

import (
	"fmt"
	"log"

	"github.com/cufee/am-stats/auth"
	"github.com/cufee/am-stats/config"
	"github.com/cufee/am-stats/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Define routes
	app := fiber.New()

	// Logger
	app.Use(logger.New())

	// API key validator
	app.Use(auth.Validator)

	// Checks
	app.Get("/player/id/:id", handlers.HandlePlayerCheck)

	// Stats
	app.Get("/player", handlers.HandleStatsImageExport)
	app.Get("/stats", handlers.HandleStatsJSONExport)

	// Achievements
	app.Get("/achievements", handlers.HandleAchievementsJSONExport)

	log.Print(app.Listen(fmt.Sprintf(":%v", config.APIport)))
}
