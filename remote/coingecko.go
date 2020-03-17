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

// CoinGeckoFactory is coinGecko's factory calss(?)
type CoinGeckoFactory struct{}

const sourceNameCoinGecko = "CoinGecko"

// Create return data after json data which form https://coinmarketcap.com/ be decoded
func (CoinGeckoFactory) Create(_authKey string) (API, error) {

	return &coinGecko{
		responseAttribute: &responseAttribute{
			sourceName: sourceNameCoinGecko,
			usd:        0.0,
			timestamp:  "",
			latestID:   primitive.ObjectID{},
			authKey:    _authKey,
		},
	}, nil
}

// GetUSD return usd about one BTC
func (cg coinGecko) CallRemote() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price", nil)
	if err != nil {
		logrus.Info(err.Error())
		return err
	}

	q := url.Values{}
	q.Add("ids", "bitcoin")
	q.Add("vs_currencies", "usd")
	q.Add("include_last_updated_at", "true")

	req.Header.Set("Accepts", "application/json")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		logrus.Infof("Error sending request to server %s", cg.sourceName)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logrus.Infof("%s return status code not equal 200", cg.sourceName)
		return err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Infof("%s body read error", cg.sourceName)
		return err
	}
	err = cg.setValues(string(respBody))
	if err != nil {
		logrus.Infof("%s json parse error", cg.sourceName)
		return err
	}
	return nil
}

func (cg *coinGecko) setValues(str string) error {
	var cgRes coinGeckoResponse
	err := json.Unmarshal([]byte(str), &cgRes)
	if err != nil {
		return err
	}
	USD := cgRes.Data
	t := time.Unix(USD.Timestamp, 0)
	if err != nil {
		logrus.Debug("time parse error")
		return err
	}
	cg.usd = USD.Price
	cg.timestamp = t.Format(time.RFC3339)

	return nil
}

type coinGecko struct {
	*responseAttribute
}

type coinGeckoData struct {
	Price     float64 `json:"usd"`
	Timestamp int64   `json:"last_updated_at"`
}

type coinGeckoResponse struct {
	Data coinGeckoData `json:"bitcoin"`
}

/* sample data
{
  "bitcoin": {
    "usd": 5310.29,
    "last_updated_at": 1584434391
  }
}
*/
