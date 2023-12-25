package requests

type GetCreateOffer struct {
	From      LatLngLiteral `json:"from"`
	To        LatLngLiteral `json:"to"`
	Client_id string        `json:"client_id"`
}

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
