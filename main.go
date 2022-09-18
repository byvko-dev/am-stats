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

	api := app.Group("/v1")

	// Public endpoints
	// Stats
	api.Post("/public/stats", handlers.HandlePublicStatsJSONExport)
	api.Post("/public/stats/reset-special", handlers.HandleSpecialSessionReset)

	// API key validator
	api.Use(auth.Validator)

	// Checks
	api.Get("/player/id/:id", handlers.HandlePlayerCheck)

	// Replays
	api.Post("/replay", handlers.HandleReplayJSONExport)
	api.Post("/replay/image", handlers.HandleReplayImageExport)

	// Stats
	api.Post("/stats", handlers.HandleStatsJSONExport)
	api.Post("/stats/image", handlers.HandleStatsImageExport)

	// Achievements
	// Clan
	api.Post("/achievements/leaderboard/clans", handlers.HandleClanAchievementsLbExport)
	api.Post("/achievements/leaderboard/clans/image", handlers.HandlerClansLeaderboardImage)
	// Players Leaderboard
	api.Post("/achievements/leaderboard/players", handlers.HandlePlayersAchievementsLbExport)
	api.Post("/achievements/leaderboard/players/image", handlers.HandlerPlayersLeaderboardImage)

	log.Panic(app.Listen(fmt.Sprintf(":%v", os.Getenv("PORT"))))
}
