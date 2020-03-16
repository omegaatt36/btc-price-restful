package remote

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoinMarketCap(t *testing.T) {
	price := 5050.98607739
	jsonStr := `{ "status": { "timestamp": "2020-03-16T06:25:18.504Z", "error_code": 0, "error_message": null, "elapsed": 10, "credit_count": 1, "notice": null }, "data": [ { "id": 1, "name": "Bitcoin", "symbol": "BTC", "slug": "bitcoin", "num_market_pairs": 7814, "date_added": "2013-04-28T00:00:00.000Z", "tags": [ "mineable" ], "max_supply": 21000000, "circulating_supply": 18272725, "total_supply": 18272725, "platform": null, "cmc_rank": 1, "last_updated": "2020-03-16T06:24:29.000Z", "quote": { "USD": { "price": ` + fmt.Sprint(price) + `, "volume_24h": 34429654396.5374, "percent_change_1h": -4.45158, "percent_change_24h": -4.4494, "percent_change_7d": -36.6381, "market_cap": 92295279570.97618, "last_updated": "2020-03-16T06:24:29.000Z" } } } ] }`
	coinFactory := CoinMarketCapFactory{}
	coin, _ := coinFactory.Create(jsonStr)
	err := coin.(*coinMarketCap).setValues(jsonStr)
	assert.Nil(t, err, "err should be nothing")
	assert.Equal(t, coin.GetUSD(), price)

	err = coin.(*coinMarketCap).setValues("test")
	assert.NotNil(t, err, "should not be nil")
}
