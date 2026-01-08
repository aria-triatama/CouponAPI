package models

type Claims struct {
	UserID     string `json:"user_id,omitempty" bson:"user_id" validate:"required"`
	CouponName string `json:"coupon_name,omitempty" bson:"coupon_name" validate:"required"`
}
