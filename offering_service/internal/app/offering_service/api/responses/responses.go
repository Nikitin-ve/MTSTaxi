package responses

type PullCreateOffer struct {
	Trip_id   string        `json:"trip_id" json:"id"`
	From      LatLngLiteral `json:"from"`
	To        LatLngLiteral `json:"to"`
	Client_id string        `json:"client_id"`
	Price     Money         `json:"price"`
}

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}
