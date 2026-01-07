package database

import (
	"CouponAPI/models"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var lookupCouponClaims = bson.D{{"$lookup", bson.D{
	{"from", "claims"},
	{"localField", "name"},
	{"foreignField", "coupon_name"},
	{"as", "claims"},
}}}

func CouponsCollection(mdb DBstore) *mongo.Collection {
	collection := mdb.GetDB().Collection("coupons")
	return collection
}

func AddCoupon(mdb DBstore, data models.Coupons) (bool, error) {
	var result models.Coupons

	coupon := CouponsCollection(mdb)

	err := coupon.FindOne(*mdb.GetCtx(), bson.M{"name": data.Name}).Decode(&result)
	if err == nil {
		return false, errors.New("coupon already exists")
	}

	_, err = coupon.InsertOne(*mdb.GetCtx(), data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CheckStock(mdb DBstore, session mongo.SessionContext, data models.Claims) bool {
	var result models.Coupons

	coupon := CouponsCollection(mdb)

	err := coupon.FindOne(session, bson.M{"name": data.CouponName}).Decode(&result)
	if err != nil {
		return false
	}

	if result.RemainingAmount <= 0 {
		return false
	}

	return true
}

func UpdateStock(mdb DBstore, session mongo.SessionContext, data models.Claims) error {
	coupon := CouponsCollection(mdb)

	_, err := coupon.UpdateOne(session, bson.M{"name": data.CouponName}, bson.M{"$inc": bson.M{"remaining_amount": -1}})
	if err != nil {
		return err
	}

	return nil
}

func GetCouponDetails(mdb DBstore, name string) (models.CouponDetails, error) {
	details := make([]models.CouponDetails, 0)

	query, err := CouponsCollection(mdb).Aggregate(*mdb.GetCtx(), mongo.Pipeline{
		lookupCouponClaims,
		bson.D{{"$unwind", bson.D{
			{"path", "$claims"},
			{"preserveNullAndEmptyArrays", true},
		}}},
		bson.D{{"$group", bson.D{
			{"_id", "$_id"},
			{"name", bson.M{"$first": "$name"}},
			{"amount", bson.M{"$first": "$amount"}},
			{"remaining_amount", bson.M{"$first": "$remaining_amount"}},
			{"claimed_by", bson.M{"$push": bson.M{"user_id": "$claims.user_id"}}},
		}}},
		bson.D{{"$match", bson.M{"name": name}}},
	})

	if err == nil {
		err = query.All(*mdb.GetCtx(), &details)
	}
	defer query.Close(*mdb.GetCtx())

	return details[0], err
}
