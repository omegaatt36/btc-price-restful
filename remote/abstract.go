package remote

import "go.mongodb.org/mongo-driver/bson/primitive"

// API is a "abstract" interface about parse remote api data for sub-class overriding
type API interface {
	GetUSD() float64
	GetSourceName() string
	GetLastestID() (primitive.ObjectID, error)
	CallRemote() error
	InsertDB() error
}

type responseAttribute struct {
	sourceName string
	usd        float64
	timestamp  string
	latestID   primitive.ObjectID
	authKey    string
}

type responseFactory interface {
	Create(string) (API, error)
}
