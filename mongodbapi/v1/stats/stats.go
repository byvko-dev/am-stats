package mongodbapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cufee/am-stats/config"
	mgo "github.com/cufee/am-stats/mongodbapi/v1"
	wgapi "github.com/cufee/am-stats/wargamingapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var sessionsCollection *mongo.Collection
var streaksCollection *mongo.Collection

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
	log.Println("Stats - successfully connected and pinged.")

	sessionsCollection = client.Database("stats").Collection("sessions")
	streaksCollection = client.Database("stats").Collection("streaks")
}

// GetPlayerSession -
func GetPlayerSession(pid int, days int, currentBattles int) (session Session, err error) {
	// Sorting options
	queryOptions := options.FindOneOptions{}
	queryOptions.SetSort(bson.M{"timestamp": -1})
	// Make a filter
	var filters []mgo.FilterPair
	filters = append(filters, mgo.FilterPair{Key: "player_id", Value: pid})
	filters = append(filters, mgo.FilterPair{Key: "battles_random", Value: bson.M{"$ne": currentBattles}})
	if days > 0 {
		// Setting days to negative to look back
		queryOptions.SetSort(bson.M{"timestamp": 1})
		sessionTime := time.Now().Add(time.Hour * 24 * -(time.Duration(days) + 1))
		filters = append(filters, mgo.FilterPair{Key: "timestamp", Value: bson.M{"$gt": sessionTime}})
	}
	query := mgo.MakeFilter(filters...)
	// Get session
	var retroSession RetroSession
	err = sessionsCollection.FindOne(ctx, query, &queryOptions).Decode(&retroSession)

	log.Print(retroSession.PlayerID)

	if err != nil {
		return session, err
	}
	// Convert to Session
	return retroSession.ToSession(), nil
}

// GetPlayerSessionAchievements -
func GetPlayerSessionAchievements(pid int, days int, fields ...string) (data wgapi.AchievementsFrame, err error) {
	// Sorting options
	queryOptions := options.FindOneOptions{}
	queryOptions.SetSort(bson.M{"timestamp": -1})

	// Generate projection
	if len(fields) > 0 {
		var project bson.D
		// Loop over field, compile project and sort
		for _, f := range fields {
			project = append(project, bson.E{Key: fmt.Sprintf("achievements.achievements.%s", f), Value: 1}) // Show field
		}
		queryOptions.SetProjection(project)
	}

	// Make a filter
	var filters []mgo.FilterPair
	filters = append(filters, mgo.FilterPair{Key: "player_id", Value: pid})
	if days > 0 {
		// Setting days to negative to look back
		queryOptions.SetSort(bson.M{"timestamp": 1})
		sessionTime := time.Now().AddDate(0, 0, -(days + 1))
		filters = append(filters, mgo.FilterPair{Key: "timestamp", Value: bson.M{"$gt": sessionTime}})
	}
	query := mgo.MakeFilter(filters...)

	// Get session
	var session Session
	err = sessionsCollection.FindOne(ctx, query, &queryOptions).Decode(&session)
	if err != nil {
		return data, err
	}
	return session.Achievements, nil
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
