package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createStringCaseInsensitiveFilter(value string) bson.M {
	return bson.M{"$regex": primitive.Regex{Pattern: value, Options: "i"}}
}
