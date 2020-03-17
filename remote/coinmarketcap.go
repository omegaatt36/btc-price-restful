package remote

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CoinMarketCapFactory is coinMarketCap's factory calss(?)
type CoinMarketCapFactory struct{}

// SourceNameCoinMarketCap is a tag let other class can direct read without creat API
const sourceNameCoinMarketCap = "CoinMarketCap"

// Create return data after json data which form https://coinmarketcap.com/ be decoded
func (CoinMarketCapFactory) Create(_authKey string) (API, error) {

	return &coinMarketCap{
		responseAttribute: &responseAttribute{
			sourceName: sourceNameCoinMarketCap,
			usd:        0.0,
			timestamp:  "",
			latestID:   primitive.ObjectID{},
			authKey:    _authKey,
		},
	}, nil
}

type coinMarketCap struct {
	*responseAttribute
}

// GetUSD return usd about one BTC
func (cmc coinMarketCap) CallRemote() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		logrus.Info(err.Error())
		return err
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "1")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", cmc.authKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		logrus.Info("Error sending request to server")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logrus.Info("coinmarketcap return status code not equal 200")
		return err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("body read error")
		return err
	}
	err = cmc.setValues(string(respBody))
	if err != nil {
		logrus.Info("json parse error")
		return err
	}
	return nil
}

func (cmc *coinMarketCap) setValues(str string) error {
	var cmcRes coinmarketcapResponse
	err := json.Unmarshal([]byte(str), &cmcRes)
	if err != nil {
		return err
	}
	USD := cmcRes.Data[0].Quote.USD
	t, err := time.Parse(time.RFC3339, USD.Timestamp)
	if err != nil {
		logrus.Debug("time parse error")
		return err
	}
	cmc.usd = USD.Price
	cmc.timestamp = t.Format("2 Jan 2006 15:04:05")

	return nil
}

type coinmarketcapUSD struct {
	Price     float64 `json:"price"`
	Timestamp string  `json:"last_updated"`
}

type coinmarketcapQuote struct {
	USD coinmarketcapUSD `json:"USD"`
}

type coinmarketcapData struct {
	Quote coinmarketcapQuote `json:"quote"`
}

type coinmarketcapResponse struct {
	Data []coinmarketcapData `json:"data"`
}

/* sample data
{
    "status": {
        "timestamp": "2020-03-16T06:25:18.504Z",
        "error_code": 0,
        "error_message": null,
        "elapsed": 10,
        "credit_count": 1,
        "notice": null
    },
    "data": [
        {
            "id": 1,
            "name": "Bitcoin",
            "symbol": "BTC",
            "slug": "bitcoin",
            "num_market_pairs": 7814,
            "date_added": "2013-04-28T00:00:00.000Z",
            "tags": [
                "mineable"
            ],
            "max_supply": 21000000,
            "circulating_supply": 18272725,
            "total_supply": 18272725,
            "platform": null,
            "cmc_rank": 1,
            "last_updated": "2020-03-16T06:24:29.000Z",
            "quote": {
                "USD": {
                    "price": 5050.98607739,
                    "volume_24h": 34429654396.5374,
                    "percent_change_1h": -4.45158,
                    "percent_change_24h": -4.4494,
                    "percent_change_7d": -36.6381,
                    "market_cap": 92295279570.97618,
                    "last_updated": "2020-03-16T06:24:29.000Z"
                }
            }
        }
    ]
}
*/
