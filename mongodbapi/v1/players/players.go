package mongodbapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cufee/am-stats/config"
	mgo "github.com/cufee/am-stats/mongodbapi/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var playersCollection *mongo.Collection

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
	log.Println("Players - successfully connected and pinged.")

	playersCollection = client.Database("stats").Collection("players")
}

// AddPlayer - Add a new player record to DB
func AddPlayer(playerData DBPlayerPofile) error {
	_, err := playersCollection.InsertOne(ctx, playerData)
	if err != nil {
		return err
	}
	return nil
}

// getPlayer - Get a player record from DB
func getPlayer(filter interface{}) (result DBPlayerPofile, err error) {
	err = playersCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetPlayerProfile - Get a player record from DB
func GetPlayerProfile(pid int) (DBPlayerPofile, error) {
	filter := mgo.MakeFilter(mgo.FilterPair{Key: "_id", Value: pid})
	return getPlayer(filter)
}

// GetRealmByPID - Get player realm by PID
func GetRealmByPID(pid int) (realm string, err error) {
	filter := mgo.MakeFilter(mgo.FilterPair{Key: "_id", Value: pid})
	profile, err := getPlayer(filter)
	return profile.Nickname, err
}

// GetRealmPlayers - Get players by realm
func GetRealmPlayers(realm string) (pidSlice []int, err error) {
	// Find
	rawSlice, err := playersCollection.Distinct(ctx, "_id", bson.M{"realm": realm})
	if err != nil {
		return pidSlice, err
	}

	// Make a slice
	for _, pid := range rawSlice {
		pidSlice = append(pidSlice, int(pid.(int32)))
	}

	return pidSlice, err
}

// UpdatePlayer - Update a player record in DB
func UpdatePlayer(filter interface{}, playerData DBPlayerPofile) (result string, err error) {
	resultRaw, err := playersCollection.UpdateOne(ctx, filter, bson.M{"$set": playerData})
	if err != nil {
		return "mongoapi/UpdatePlayer: Error updating player record.", err
	}
	result = fmt.Sprintf("%+v", resultRaw)
	return result, nil
}
