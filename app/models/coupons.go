package models

type Coupons struct {
	Name            string `json:"name,omitempty" bson:"name"`
	Amount          int    `json:"amount,omitempty" bson:"amount"`
	RemainingAmount int    `json:"remaining_amount,omitempty" bson:"remaining_amount"`
}

type CouponsRequest struct {
	Name   string `json:"name,omitempty" validate:"required"`
	Amount int    `json:"amount,omitempty" validate:"required"`
}
