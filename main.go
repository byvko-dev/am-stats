package main

import (
	"fmt"
	"log"

	"github.com/cufee/am-stats/auth"
	"github.com/cufee/am-stats/config"
	"github.com/cufee/am-stats/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Define routes
	app := fiber.New()

	// Logger
	app.Use(logger.New())
	// CORS
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Public endpoints
	// Stats
	app.Post("/public/stats", handlers.HandlePublicStatsJSONExport)
	app.Post("/public/stats/reset-special", handlers.HandleSpecialSessionReset)

	// API key validator
	app.Use(auth.Validator)

	// Checks
	app.Get("/player/id/:id", handlers.HandlePlayerCheck)

	// Replays
	app.Get("/replay", handlers.HandleReplayJSONExport)
	app.Get("/replay/image", handlers.HandleReplayImageExport)

	// Stats
	app.Get("/stats", handlers.HandleStatsJSONExport)
	app.Get("/stats/image", handlers.HandleStatsImageExport)

	// Achievements
	// Clan
	app.Get("/achievements/leaderboard/clans", handlers.HandleClanAchievementsLbExport)
	app.Get("/achievements/leaderboard/clans/image", handlers.HandlerClansLeaderboardImage)
	// Players Leaderboard
	app.Get("/achievements/leaderboard/players", handlers.HandlePlayersAchievementsLbExport)
	app.Get("/achievements/leaderboard/players/image", handlers.HandlerPlayersLeaderboardImage)

	log.Print(app.Listen(fmt.Sprintf(":%v", config.APIport)))
}
