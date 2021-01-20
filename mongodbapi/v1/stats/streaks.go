package mongodbapi

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
