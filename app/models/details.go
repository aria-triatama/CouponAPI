package models

type CouponDetails struct {
	Name            string   `json:"name,omitempty" bson:"name"`
	Amount          int      `json:"amount,omitempty" bson:"amount"`
	RemainingAmount int      `json:"remaining_amount,omitempty" bson:"remaining_amount"`
	ClaimedBy       []string `json:"claimed_by,omitempty" bson:"claimed_by"`
}
