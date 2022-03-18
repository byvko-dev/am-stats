package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cufee/am-stats/auth"
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

	api := app.Group("/stats/v1")

	// Public endpoints
	// Stats
	api.Post("/public/stats", handlers.HandlePublicStatsJSONExport)
	api.Post("/public/stats/reset-special", handlers.HandleSpecialSessionReset)

	// API key validator
	api.Use(auth.Validator)

	// Checks
	api.Get("/player/id/:id", handlers.HandlePlayerCheck)

	// Replays
	api.Get("/replay", handlers.HandleReplayJSONExport)
	api.Get("/replay/image", handlers.HandleReplayImageExport)

	// Stats
	api.Get("/stats", handlers.HandleStatsJSONExport) // Legacy
	api.Post("/stats", handlers.HandleStatsJSONExport)
	api.Get("/stats/image", handlers.HandleStatsImageExport)

	// Achievements
	// Clan
	api.Get("/achievements/leaderboard/clans", handlers.HandleClanAchievementsLbExport)
	api.Get("/achievements/leaderboard/clans/image", handlers.HandlerClansLeaderboardImage)
	// Players Leaderboard
	api.Get("/achievements/leaderboard/players", handlers.HandlePlayersAchievementsLbExport)
	api.Get("/achievements/leaderboard/players/image", handlers.HandlerPlayersLeaderboardImage)

	log.Panic(app.Listen(fmt.Sprintf(":%v", os.Getenv("PORT"))))
}
