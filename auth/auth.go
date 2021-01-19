package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// DefaultHeaderKeyIdentifier is the default api key identifier in request headers
	DefaultHeaderKeyIdentifier string = "x-api-key"
)

// Validator - Validate API key passed in header
func Validator(c *fiber.Ctx) error {
	// Parse api key
	headerKey := c.Get(DefaultHeaderKeyIdentifier)

	// Check if API key was provided
	if headerKey == "" {
		return fiber.ErrBadRequest
	}

	// Get app data
	appData, valid := validateKey(headerKey)

	// Check if the key is enabled
	if valid {
		defer func() {
			// Generate IP warning
			if appData.LastIP != c.IP() {
				log.Print(fmt.Sprintf("Application %s changed IP address from %s to %s", appData.AppName, appData.LastIP, c.IP()))

				// Update last used IP
				go updateAppLastIP(appData.AppID, c.IP())
			}

			// Log request
			go logEvent(appData, *c)
		}()

		// Go to next middleware:
		return c.Next()
	}
	return fiber.ErrUnauthorized
}

// validateKey - Validate API Key
func validateKey(key string) (appData appllicationData, valid bool) {
	// Get application info
	appData, err := appDataByKey(key)

	// Log error
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Print(fmt.Errorf("appDataByKey: %s", err.Error()))
	}

	// Return
	return appData, appData.Enabled
}

// updateAppLastIP - Update last IP used for app
func updateAppLastIP(appID primitive.ObjectID, IP string) {
	var appData appllicationData
	appData.AppID = appID
	appData.LastIP = IP
	appData.LastUsed = time.Now()

	err := updateAppData(appData)
	if err != nil {
		log.Print(fmt.Errorf("updateAppData: %s", err.Error()))
	}
}

// logEvent - Log access event
func logEvent(appData appllicationData, c fiber.Ctx) {
	// Prepare log data
	logData, err := appData.prepLogData()
	if err != nil {
		log.Print(fmt.Errorf("prepLogData: %s", err.Error()))
		return
	}

	// Fill log data
	logData.RequestIP = c.IP()
	logData.RequestPath = c.Path()
	logData.RequestTime = time.Now()
	logData.RequestMethod = c.Method()

	err = addLogEntry(logData)
	if err != nil {
		log.Print(fmt.Errorf("addLogEntry: %s", err.Error()))
	}
}
