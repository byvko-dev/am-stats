package mongodbapi

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/cufee/am-stats/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var tankAveragesCollection *mongo.Collection
var tankGlossaryCollection *mongo.Collection
var achievementsGlossaryCollection *mongo.Collection

// Ctx - Context for MongoDB connection
var ctx context.Context

// Client - Client for MongoDB connection
var client *mongo.Client

func init() {
	// Conenct to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Println("Panic in mongoapi/init")
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Panic in mongoapi/init")
		panic(err)
	}
	log.Println("Glossary - successfully connected and pinged.")

	tankAveragesCollection = client.Database("glossary").Collection("tankaverages")
	tankGlossaryCollection = client.Database("glossary").Collection("tanks")
	achievementsGlossaryCollection = client.Database("glossary").Collection("achievements")
}

// GetTankAverages - Get averages data for a tank by ID
func GetTankAverages(tid int) (averages TankAverages, err error) {
	filter := bson.M{"tank_id": tid}
	err = tankAveragesCollection.FindOne(ctx, filter).Decode(&averages)
	return averages, err
}

// GetAchievementIcon - Get an icon for achievement
func GetAchievementIcon(aid string) (url string, err error) {
	filter := bson.M{"_id": strings.ToLower(aid)}
	result, err := achievementsGlossaryCollection.Distinct(ctx, "image", filter)
	if len(result) > 0 {
		return result[0].(string), err
	}
	return url, err
}

// GetTankGlossary - Get averages data for a tank by ID
func GetTankGlossary(tid int) (averages TankAverages, err error) {
	filter := bson.M{"tank_id": tid}
	err = tankGlossaryCollection.FindOne(ctx, filter).Decode(&averages)
	return averages, err
}
