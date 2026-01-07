package models

type Claims struct {
	UserID     string `json:"user_id,omitempty" validate:"required"`
	CouponName string `json:"coupon_name,omitempty" validate:"required"`
}
