package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User provides a regular struct to get JWT, and it will save in db.
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password,omitempty"`
}

// JwtToken provides authority to access service
type JwtToken struct {
	Token string `json:"token"`
}
