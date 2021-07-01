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
var achievementsCacheCollection *mongo.Collection

// Ctx - Context for MongoDB connection
var ctx context.Context

func init() {
	// Connect to MongoDB
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

	achievementsClansCollection = client.Database("achievements").Collection("clans")
	achievementsCacheCollection = client.Database("achievements").Collection("cache")
	achievementsPlayersCollection = client.Database("achievements").Collection("players")
}

// Check cached request
func CheckCachedMedals(realm string, days int, medals []MedalWeight, expiration time.Duration) (export []AchievementsPlayerData, totalScore int, err error) {
	// Find cached request
	var cache CachedMedalsRequest
	var filter bson.M = bson.M{"request.realm": realm, "request.days": days, "request.medals": medals, "result.updated_timestamp": bson.M{"$gt": time.Now().Add(-expiration)}}
	err = achievementsCacheCollection.FindOneAndUpdate(ctx, filter, bson.M{"$set": bson.M{"requested_timestamp": time.Now()}}).Decode(&cache)
	return cache.Result.SortedPlayers, cache.Result.TotalScore, err
}

// Check cached request
func SaveCachedMedals(realm string, days int, medals []MedalWeight, sortedPlayers []AchievementsPlayerData, totalScore int) {
	opts := options.FindOneAndReplace()
	opts.SetUpsert(true)
	var cache CachedMedalsRequest
	cache.Request.Days = days
	cache.Request.Realm = realm
	cache.Request.Medals = medals
	cache.LastRequested = time.Now()
	cache.Result.UpdatedAt = time.Now()
	cache.Result.TotalScore = totalScore
	cache.Result.SortedPlayers = sortedPlayers
	achievementsCacheCollection.FindOneAndReplace(ctx, bson.M{"request": cache.Request}, cache, opts)
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

// GetPlayerAchievements - Get last cached player achievements from player ID
func GetPlayerAchievements(pid int, medals ...MedalWeight) (data AchievementsPlayerData, err error) {
	opts := options.FindOne()
	// Generate projection
	if len(medals) > 0 {
		var project bson.D
		project = append(project, bson.E{Key: "timestamp", Value: 1}) // Show timestamp
		// Loop over field, compile project and sort
		for _, m := range medals {
			query := fmt.Sprintf("data.achievements.%s", m.Name)
			project = append(project, bson.E{Key: query, Value: 1}) // Show field
		}
		opts.Projection = project
	}

	// Find Player
	err = achievementsPlayersCollection.FindOne(ctx, bson.M{"_id": pid}, opts).Decode(&data)
	if err != nil {
		return
	}

	return data, err
}

// GetClanAchievementsCache - Get last cached clans achievements leaderboard
func GetClanAchievementsCache(clanID int, fields ...string) (data []AchievementsPlayerData, err error) {
	opts := options.Find()
	// Generate projection
	if len(fields) > 0 {
		var project bson.D
		// Loop over field, compile project and sort
		for _, f := range fields {
			project = append(project, bson.E{Key: fmt.Sprintf("data.achievements.%s", f), Value: 1}) // Show field
		}
		project = append(project, bson.E{Key: "_id", Value: 1})      // Always show Clan ID
		project = append(project, bson.E{Key: "clan_tag", Value: 1}) // Always show clan tag
		project = append(project, bson.E{Key: "members", Value: 1})  // Always show members
		opts.Projection = project
	}

	// Find
	cur, err := achievementsClansCollection.Find(ctx, bson.M{"_id": clanID}, opts)
	if err != nil {
		return data, err
	}

	// Decode and return
	if err = cur.All(ctx, &data); err != nil {
		return data, err
	}
	return data, err
}

// SearchClanAchievementsLb - Get last cached achievements data for clan tag and realm
func SearchClanAchievementsLb(tag string, realm string) (data ClanAchievements, err error) {
	filter := bson.M{"clan_tag": tag, "realm": realm}
	err = achievementsClansCollection.FindOne(ctx, filter).Decode(&data)
	return data, err
}
