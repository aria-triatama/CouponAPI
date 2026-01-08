package database

import (
	"CouponAPI/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ClaimsCollection(mdb DBstore) *mongo.Collection {
	collection := mdb.GetDB().Collection("claims")
	return collection
}

func EnsureClaimsIndexes(mdb DBstore) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "coupon_name", Value: 1},
		},
		Options: options.Index().
			SetUnique(true),
	}

	claim := ClaimsCollection(mdb)
	_, err := claim.Indexes().CreateOne(*mdb.GetCtx(), indexModel)
	return err
}

func AddClaim(mdb DBstore, data models.Claims) error {
	claim := ClaimsCollection(mdb)

	_, err := claim.InsertOne(*mdb.GetCtx(), data)
	if err != nil {
		return err
	}

	return nil
}

func CheckClaim(mdb DBstore, data models.Claims) bool {
	var result models.Coupons

	claim := ClaimsCollection(mdb)

	err := claim.FindOne(*mdb.GetCtx(), bson.D{{"user_id", data.UserID}, {"coupon_name", data.CouponName}}).Decode(&result)
	if err != nil {
		return false
	}

	return true
}
