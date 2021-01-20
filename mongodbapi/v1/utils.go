package mongodbapi

import (
	"go.mongodb.org/mongo-driver/bson"
)

// FilterPair - Used to make BSON filters
type FilterPair struct {
	Key   string
	Value interface{}
}

// MakeFilter - Make a BSON filter for mongodb using FilerPairs passed in
func MakeFilter(filters ...FilterPair) (filter interface{}) {
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
