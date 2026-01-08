package models

import "time"

type Gates struct {
	ID        string    `bson:"id"`
	CreatedAt time.Time `bson:"created_at"`
}
