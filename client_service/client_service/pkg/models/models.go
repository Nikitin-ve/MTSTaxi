package models

type Trip struct {
	User_id  string        `json:"user_id"`
	Offer_id string        `json:"offer_id"`
	Trip_id  string        `json:"trip_id" json:"id"`
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
	Price    Money         `json:"price"`
	Status   string        `json:"status"`
}

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}
