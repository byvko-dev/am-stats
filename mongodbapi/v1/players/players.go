package mongodbapi

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/cufee/am-stats/config"
	mgo "github.com/cufee/am-stats/mongodbapi/v1"
	wgapi "github.com/cufee/am-stats/wargamingapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var playersCollection *mongo.Collection
var pinsCollection *mongo.Collection

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
	log.Println("Players - successfully connected and pinged.")

	playersCollection = client.Database("stats").Collection("players")
	pinsCollection = client.Database("stats").Collection("pins")
}

// GetPlayerProfile - Get a player record from DB
func GetPinsBulk(pinIDs ...primitive.ObjectID) ([]UserPin, error) {
	var pins []UserPin
	if len(pinIDs) == 0 {
		return pins, nil
	}
	filter := bson.M{"_id": bson.M{"$in": pinIDs}, "hidden": bson.M{"$ne": true}}
	cur, err := pinsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return pins, cur.All(ctx, &pins)
}

// AddPlayer - Add a new player record to DB
func AddPlayer(playerData DBPlayerProfile) error {
	_, err := playersCollection.InsertOne(ctx, playerData)
	if err != nil {
		return err
	}
	return nil
}

// getPlayer - Get a player record from DB
func getPlayer(filter interface{}) (result DBPlayerProfile, err error) {
	err = playersCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetPlayerProfile - Get a player record from DB
func GetPlayerProfile(pid int) (DBPlayerProfile, error) {
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
		if reflect.TypeOf(pid) == reflect.TypeOf(int32(1)) {
			pidSlice = append(pidSlice, int(pid.(int32)))
		}
	}

	return pidSlice, err
}

// UpdatePlayer - Update a player record in DB
func UpdatePlayer(filter interface{}, playerData DBPlayerProfile) (result string, err error) {
	resultRaw, err := playersCollection.UpdateOne(ctx, filter, bson.M{"$set": playerData})
	if err != nil {
		return "mongoapi/UpdatePlayer: Error updating player record.", err
	}
	result = fmt.Sprintf("%+v", resultRaw)
	return result, nil
}

// CheckUserByUserID - Check user profile by Discord ID
func GreedyClanPlayerCapture(player wgapi.PlayerProfile, realm string) {
	// Get full clan profile
	clan, err := wgapi.ClanDataByID(player.ClanID, realm)
	if err != nil {
		log.Print("Failed to greedy capture clan data ", err.Error())
		return
	}

	// Add new players to DB
	for _, m := range clan.MembersIds {
		profile, err := GetPlayerProfile(m)
		if profile.ID != 0 || (err != nil && err.Error() != "mongo: no documents in result") {
			if err != nil {
				log.Print("Error during greedy capture GetPlayerProfile - ", err.Error())
			}
			continue
		}
		profile.ID = m
		profile.CareerWN8 = -1
		profile.ClanID = player.ClanID
		profile.ClanTag = player.ClanTag
		profile.ClanName = player.ClanName
		profile.Realm = strings.ToUpper(realm)
		err = AddPlayer(profile)
		if err != nil {
			log.Print("Error during greedy capture AddPlayer - ", err.Error())
		}
	}
}
