package requests

type GetUserID struct {
	User_id string `json:"user_id"`
}

type ForDB struct {
	Trip_id         string  `json:"id"`
	Source          string  `json:"source"`
	Type            string  `json:"type"`
	Datacontenttype string  `json:"datacontenttype"`
	Time            string  `json:"time"`
	Offer_id        string  `json:"offer_id"`
	Amount          int64   `json:"amount"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	Lat_from        float64 `json:"lat"`
	Lng_from        float64 `json:"lng"`
	Lat_to          float64 `json:"lat"`
	Lng_to          float64 `json:"lng"`
}

type CreatedTrip struct {
	Trip_id         string `json:"id"`
	Source          string `json:"source"`
	Type            string `json:"type"`
	Datacontenttype string `json:"datacontenttype"`
	Time            string `json:"time"`
	Data            Offer5 `json:"data"`
}

type Offer5 struct {
	Trip_id  string        `json:"trip_id"`
	Offer_id string        `json:"offer_id"`
	Price    Money         `json:"price"`
	Status   string        `json:"status"`
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
}

type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
