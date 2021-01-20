package mongodbapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cufee/am-stats/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var achievementsPlayersCollection *mongo.Collection
var achievementsClansCollection *mongo.Collection

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
	log.Println("Achievements - successfully connected and pinged.")

	achievementsPlayersCollection = client.Database("achievements").Collection("players")
	achievementsClansCollection = client.Database("achievements").Collection("clans")
}

// GetPlayerAchievementsLb - Get last cached players achievements leaderboard
func GetPlayerAchievementsLb(realm string, fields ...string) (data []AchievementsPlayerData, err error) {
	opts := options.Find()
	// Generate projection
	if len(fields) > 0 {
		var project bson.D
		// Loop over field, compile project and sort
		for _, f := range fields {
			project = append(project, bson.E{Key: fmt.Sprintf("data.achievements.%s", f), Value: 1}) // Show field
		}
		project = append(project, bson.E{Key: "_id", Value: 1}) // Always show player ID
		opts.Projection = project
	}

	// Find
	cur, err := achievementsPlayersCollection.Find(ctx, bson.M{"realm": realm}, opts)
	if err != nil {
		return data, err
	}

	// Decode and return
	if err = cur.All(ctx, &data); err != nil {
		return data, err
	}
	return data, err
}

// GetPlayerAchievementsCache - Get last cached achievements data for player by ID
func GetPlayerAchievementsCache(pid int) (data AchievementsPlayerData, err error) {
	filter := bson.M{"_id": pid}
	err = achievementsPlayersCollection.FindOne(ctx, filter).Decode(&data)
	return data, err
}

// GetClanAchievementsCache - Get last cached achievements data for clan by ID
func GetClanAchievementsCache(clanID int) (data ClanAchievements, err error) {
	filter := bson.M{"_id": clanID}
	err = achievementsClansCollection.FindOne(ctx, filter).Decode(&data)
	return data, err
}

// SearchClanAchievementsCache - Get last cached achievements data for clan tag and realm
func SearchClanAchievementsCache(tag string, realm string) (data ClanAchievements, err error) {
	filter := bson.M{"clan_tag": tag, "realm": realm}
	err = achievementsClansCollection.FindOne(ctx, filter).Decode(&data)
	return data, err
}
