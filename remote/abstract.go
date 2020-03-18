package remote

import (
	"BTC-price-restful/db"
	"BTC-price-restful/models"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// API is a "abstract" interface about parse remote api data for sub-class overriding
type API interface {
	GetUSD() float64
	GetSourceName() string
	GetLastestID() (primitive.ObjectID, error)
	CallRemote() error
	InsertDB() error
	InitFormRedis() error
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

// InitFormRedis is second way to init API when first call error.
func (api responseAttribute) InitFormRedis() error {
	hexID, err := db.RedisGet(fmt.Sprintf("latestAPI:%s", api.sourceName))
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}
	api.latestID = id
	return nil
}

func (api responseAttribute) InsertDB() error {
	name := api.sourceName
	var priceInfo models.Price
	priceInfo.ID = primitive.NewObjectID()
	priceInfo.Price = api.usd
	priceInfo.Source = name
	priceInfo.Timestamp = api.timestamp
	_, err := db.Create(db.CollectionPrice, priceInfo)
	if err != nil {
		logrus.Infof("%s db insert latest price info error", name)
		return err
	}
	b, _ := json.Marshal(priceInfo)
	db.RedisSet(fmt.Sprintf("latestAPI:%s", name), string(b))
	api.latestID = priceInfo.ID
	return nil
}

func (api responseAttribute) GetLastestID() (primitive.ObjectID, error) {
	z := primitive.ObjectID{}
	if api.latestID == z {
		return api.latestID, errors.New("not yet")
	}
	return api.latestID, nil
}

// GetUSD return usd about one BTC
func (api responseAttribute) GetUSD() float64 {
	return api.usd
}

// GetSourceName return this API sourceName
func (api responseAttribute) GetSourceName() string {
	return api.sourceName
}
