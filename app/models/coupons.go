package models

type Coupons struct {
	Name            string `json:"name,omitempty"`
	Amount          int    `json:"amount,omitempty"`
	RemainingAmount int    `json:"remaining_amount,omitempty"`
}

type CouponsRequest struct {
	Name   string `json:"name,omitempty" validate:"required"`
	Amount int    `json:"amount,omitempty" validate:"required"`
}
