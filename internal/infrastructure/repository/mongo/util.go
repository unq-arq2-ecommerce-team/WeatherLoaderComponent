package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	matchOp = "$match"
	groupOp = "$group"
	avgOp   = "$avg"
	gteOp   = "$gte"
	lteOp   = "$lte"
)

func createStringCaseInsensitiveFilter(value string) bson.M {
	return bson.M{"$regex": primitive.Regex{Pattern: fmt.Sprintf("^%s", value), Options: "i"}}
}

func getAverageTemperatureByCityAndDatesPipeline(city string, dateFrom, dateTo time.Time) mongo.Pipeline {
	matchStage := bson.D{
		{matchOp, bson.D{
			{"city", city},
			{"timestamp", bson.D{
				{gteOp, dateFrom},
				{lteOp, dateTo},
			}},
		}},
	}
	groupStage := bson.D{
		{groupOp, bson.D{
			{"_id", "$city"},
			{"avgTemperature", bson.D{
				{avgOp, bson.D{
					{avgOp, bson.A{"$temperatureMax", "$temperatureMin"}}},
				},
			}},
		}},
	}
	return mongo.Pipeline{
		matchStage,
		groupStage,
	}
}
