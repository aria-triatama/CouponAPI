package database

import (
	"CouponAPI/models"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
