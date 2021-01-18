package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cufee/am-stats/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Application data type
type appllicationData struct {
	AppID    primitive.ObjectID `bson:"_id,omitempty"`
	APIKey   string             `bson:"api_key,omitempty"`
	Enabled  bool               `bson:"key_enabled,omitempty"`
	AppName  string             `bson:"app_name,omitempty"`
	LastIP   string             `bson:"last_ip,omitempty"`
	LastUsed time.Time          `bson:"last_used,omitempty"`
}

// Application log entry
type appLogEntry struct {
	AppID       primitive.ObjectID `bson:"app_id"`
	AppName     string             `bson:"app_name"`
	RequestIP   string             `bson:"request_ip"`
	RequestPath string             `bson:"request_path"`
	RequestTime time.Time          `bson:"request_time"`
}

func (appData appllicationData) prepLogData() (logData appLogEntry, err error) {
	// Check application ID
	if appData.AppID == (primitive.ObjectID{}) {
		return logData, fmt.Errorf("invalid application id")
	}

	logData.AppID = appData.AppID
	logData.AppName = appData.AppName
	return logData, err
}

// Collections
var authApplicationCollection *mongo.Collection
var authLogsCollection *mongo.Collection
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
	log.Println("Auth successfully connected and pinged.")

	// Collections
	authApplicationCollection = client.Database("auth").Collection("applications")
	authLogsCollection = client.Database("auth").Collection("logs")
}

// appDataByKey - Get application info from API key
func appDataByKey(key string) (appData appllicationData, err error) {
	err = authApplicationCollection.FindOne(ctx, bson.M{"api_key": key}).Decode(&appData)
	return appData, err
}

// appDataByKey - Get application info from API key
func updateAppData(appData appllicationData) error {
	_, err := authApplicationCollection.UpdateOne(ctx, bson.M{"_id": appData.AppID}, bson.M{"$set": appData})
	return err
}

// addLogEntry - Add log entry for application
func addLogEntry(entry appLogEntry) error {
	_, err := authLogsCollection.InsertOne(ctx, entry)
	return err
}
