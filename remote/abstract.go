package remote

// Response is a "abstract" interface about parse remote api data for sub-class overriding
type Response interface {
	GetUSD() float64
}

type responseAttribute struct {
	usd float64
}

type responseFactory interface {
	Create(string) (Response, error)
}
