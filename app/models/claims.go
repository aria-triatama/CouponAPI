package models

type Claims struct {
	UserID     string `json:"user_id,omitempty"`
	CouponName string `json:"coupon_name,omitempty"`
}
