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

func EnsureIndexes(mdb DBstore) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"user_id", "1"},
			{"coupon_name", "1"},
		},
		Options: options.Index().
			SetUnique(true),
	}

	claim := ClaimsCollection(mdb)
	_, err := claim.Indexes().CreateOne(*mdb.GetCtx(), indexModel)
	return err
}

func AddClaim(mdb DBstore, session mongo.SessionContext, data models.Claims) error {
	claim := ClaimsCollection(mdb)

	_, err := claim.InsertOne(session, data)
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
