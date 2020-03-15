package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

// SetClint initialize client
func SetClint(c *mongo.Client) {
	// can be more explicitly
	if client == nil {
		client = c // <--- NOT THREAD SAFE
	}
}
