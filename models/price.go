package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Price provides a regular struct to get price, and it will save in db.
type Price struct {
	ID        primitive.ObjectID `json:"-" bson:"_id"`
	Source    string             `json:"source" bson:"source"`
	Price     float64            `json:"price" bson:"price,omitempty"`
	Timestamp string             `json:"timestamp" bson:"timestamp,omitempty"`
}
