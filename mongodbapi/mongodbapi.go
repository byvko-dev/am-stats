package mongodbapi

import (
	"fmt"
	"log"
	"time"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/cufee/am-stats/config"
)

// Collections
var sessionsCollection *mongo.Collection
var playersCollection *mongo.Collection
var streaksCollection *mongo.Collection
var tankAveragesCollection *mongo.Collection
var tankGlossaryCollection *mongo.Collection
var ctx = context.TODO()

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
	log.Println("Successfully connected and pinged.")

	// Collections
	sessionsCollection = client.Database("stats").Collection("sessions")
	playersCollection = client.Database("stats").Collection("players")
	streaksCollection = client.Database("stats").Collection("streaks")
	tankAveragesCollection = client.Database("glossary").Collection("tankaverages")
	tankGlossaryCollection = client.Database("glossary").Collection("tanks")
}

// MakeFilter - Make a BSON filter for mongodb using FilerPairs passed in
func makeFilter(filters ...FilterPair) (filter interface{}) {
	if len(filters) == 1 {
		filter = bson.M{filters[0].Key: filters[0].Value}
	}
	if len(filters) > 1 {
		var query []bson.M
		for _, v := range filters {
			query = append(query, bson.M{v.Key: v.Value})
		}
		filter = bson.M{"$and": query}
	}
	return filter
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
	filter := makeFilter(FilterPair{Key: "_id", Value: pid})
	return getPlayer(filter)
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

// GetPlayerSession -
func GetPlayerSession(pid int, days int, currentBattles int) (session Session, err error) {
	// Sorting options
	queryOptions := options.FindOneOptions{}
	queryOptions.SetSort(bson.M{"timestamp": -1})
	// Make a filter
	var filters []FilterPair
	filters = append(filters, FilterPair{Key: "player_id", Value: pid})
	filters = append(filters, FilterPair{Key: "battles_random", Value: bson.M{"$ne": currentBattles}})
	if days > 0 {
		// Setting days to negative to look back
		queryOptions.SetSort(bson.M{"timestamp": 1})
		sessionTime := time.Now().AddDate(0, 0, -(days + 1))
		filters = append(filters, FilterPair{Key: "timestamp", Value: bson.M{"$gt": sessionTime}})
	}
	query := makeFilter(filters...)
	// Get session
	var retroSession RetroSession
	err = sessionsCollection.FindOne(ctx, query, &queryOptions).Decode(&retroSession)
	if err != nil {
		return session, err
	}
	// Comvert to Session
	var retroConv Convert = retroSession
	return retroConv.ToSession(), nil
}

// AddSession - Add a new session to db
func AddSession(session Session) error {
	// Timestamp
	session.Timestamp = time.Now()
	// Concert to RetroSession
	var sessionConv Convert = session
	_, err := sessionsCollection.InsertOne(ctx, sessionConv.ToRetro())
	return err
}

// GetSession - Get a player session from db using advanced BSON filter
func GetSession(filter interface{}) (session Session, err error) {
	var retroSession RetroSession
	err = sessionsCollection.FindOne(ctx, filter).Decode(&retroSession)
	if err != nil {
		return session, err
	}
	// Comvert to Session
	var retroConv Convert = retroSession
	return retroConv.ToSession(), nil
}

// GetTankAverages - Get averages data for a tank by ID
func GetTankAverages(tid int) (averages TankAverages, err error) {
	filter := bson.M{"tank_id": tid}
	err = tankAveragesCollection.FindOne(ctx, filter).Decode(&averages)
	return averages, err
}

// GetTankGlossary - Get averages data for a tank by ID
func GetTankGlossary(tid int) (averages TankAverages, err error) {
	filter := bson.M{"tank_id": tid}
	err = tankGlossaryCollection.FindOne(ctx, filter).Decode(&averages)
	return averages, err
}

// GetStreak - Get win streak for a player by playerID
func GetStreak(pid int) (streak PlayerStreak, err error) {
	filter := bson.M{"_id": pid}
	err = streaksCollection.FindOne(ctx, filter).Decode(&streak)
	return streak, err
}

// UpdateStreak - Update win streak for a player by playerID
func UpdateStreak(streak PlayerStreak) (err error) {
	if streak.PlayerID == nil || streak.Battles == nil || streak.Losses == nil {
		// Streak data is incomplete
		return fmt.Errorf("invalid streak data passed in")
	}
	filter := bson.M{"_id": streak.PlayerID}
	streak.Timestamp = time.Now()
	update := bson.M{"$set": streak}
	opts := options.Update().SetUpsert(true)
	_, err = streaksCollection.UpdateOne(ctx, filter, update, opts)
	return err
}
