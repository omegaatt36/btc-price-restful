package remote

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CoinDeskFactory is coinDesk's factory calss(?)
type CoinDeskFactory struct{}

const sourceNameCoinDesk = "CoinDesk"

// Create return data after json data which form https://coinmarketcap.com/ be decoded
func (CoinDeskFactory) Create(_authKey string) (API, error) {

	return &coinDesk{
		responseAttribute: &responseAttribute{
			sourceName: sourceNameCoinDesk,
			usd:        0.0,
			timestamp:  "",
			latestID:   primitive.ObjectID{},
			authKey:    _authKey,
		},
	}, nil
}

// GetUSD return usd about one BTC
func (cd coinDesk) CallRemote() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coindesk.com/v1/bpi/currentprice/USD.json", nil)
	if err != nil {
		logrus.Info(err.Error())
		return err
	}

	req.Header.Set("Accepts", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Infof("Error sending request to server %s", cd.sourceName)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logrus.Infof("%s return status code not equal 200", cd.sourceName)
		return err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Infof("%s body read error", cd.sourceName)
		return err
	}
	err = cd.setValues(string(respBody))
	if err != nil {
		logrus.Infof("%s json parse error", cd.sourceName)
		return err
	}
	return nil
}

func (cd *coinDesk) setValues(str string) error {
	var cgRes coinDeskResponse
	err := json.Unmarshal([]byte(str), &cgRes)
	if err != nil {
		return err
	}
	cd.usd = cgRes.BPI.USD.Price
	cd.timestamp = cgRes.Time.Timestamp

	return nil
}

type coinDesk struct {
	*responseAttribute
}

type coinDeskUSD struct {
	Price float64 `json:"rate_float"`
}
type coinDeskBPI struct {
	USD coinDeskUSD `json:"USD"`
}
type coinDeskTime struct {
	Timestamp string `json:"updatedISO"`
}

type coinDeskResponse struct {
	Time coinDeskTime `json:"time"`
	BPI  coinDeskBPI  `json:"bpi"`
}

/* sample data
{"time":{"updated":"Mar 17, 2020 13:00:00 UTC","updatedISO":"2020-03-17T13:00:00+00:00","updateduk":"Mar 17, 2020 at 13:00 GMT"},"disclaimer":"This data was produced from the CoinDesk Bitcoin Price Index (USD). Non-USD currency data converted using hourly conversion rate from openexchangerates.org","bpi":{"USD":{"code":"USD","rate":"5,228.9150","description":"United States Dollar","rate_float":5228.915}}}
*/
