package remote

import (
	"encoding/json"
)

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

// CoinMarketCapFactory is coinMarketCap's factory calss(?)
type CoinMarketCapFactory struct{}

// Create return data after json data which form https://coinmarketcap.com/ be decoded
func (CoinMarketCapFactory) Create(str string) (Response, error) {
	var cmc coinmarketcapResponse
	err := json.Unmarshal([]byte(str), &cmc)
	if err != nil {
		return nil, err
	}
	return &coinMarketCap{
		responseAttribute: &responseAttribute{
			usd: cmc.Data[0].Quote.USD.Price,
		},
	}, nil
}

type coinMarketCap struct {
	*responseAttribute
}

// GetUSD return usd about one BTC
func (cmc coinMarketCap) GetUSD() float64 {
	return cmc.usd
}

type coinmarketcapUSD struct {
	Price float64 `json:"price"`
}

type coinmarketcapQuote struct {
	USD coinmarketcapUSD `json:"USD"`
}

type coinmarketcapData struct {
	Quote coinmarketcapQuote `json:"quote"`
}

// type coinmarketcapStatus struct {
// 	Timestamp string `json:"timestamp"`
// }

type coinmarketcapResponse struct {
	// status coinmarketcapStatus `json:"status"`
	Data []coinmarketcapData `json:"data"`
}
