package models

type CouponDetails struct {
	Name            string   `json:"name,omitempty"`
	Amount          int      `json:"amount,omitempty"`
	RemainingAmount int      `json:"remaining_amount,omitempty"`
	ClaimedBy       []string `json:"claimed_by,omitempty"`
}
