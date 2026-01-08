package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GatesCollection(mdb DBstore) *mongo.Collection {
	collection := mdb.GetDB().Collection("gates")
	return collection
}

func EnsureGatesIndexes(mdb DBstore) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "id", Value: 1},
		},
		Options: options.Index().
			SetUnique(true),
	}

	gates := GatesCollection(mdb)
	_, err := gates.Indexes().CreateOne(*mdb.GetCtx(), indexModel)
	return err
}

func ReserveGate(mdb DBstore, id string) bool {
	gates := GatesCollection(mdb)

	res, err := gates.UpdateOne(
		*mdb.GetCtx(),
		bson.M{"id": id},
		bson.M{
			"$setOnInsert": bson.M{
				"created_at": time.Now(),
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return false
	}

	if res.UpsertedCount == 0 {
		return false
	}

	return true
}

func ClearGate(mdb DBstore, id string) bool {
	gates := GatesCollection(mdb)

	_, err := gates.DeleteOne(*mdb.GetCtx(), bson.M{"id": id})
	if err != nil {
		return false
	}

	return true
}
