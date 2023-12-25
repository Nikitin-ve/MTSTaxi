package handlers

import (
	"encoding/json"
	"github.com/AlexanderGrom/componenta/crypt"
	"github.com/gorilla/mux"
	"io"
	"math"
	"net/http"
	"offering_service/internal/app/offering_service/api/requests"
	"offering_service/internal/app/offering_service/api/responses"
)

const key = "jkaheg5w8gyAKFG92"

type TaxiHandler interface {
	CreateOffer(w http.ResponseWriter, r *http.Request)
	ParseOffer(w http.ResponseWriter, r *http.Request)
}

func CreateOffer(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Problem for request", http.StatusInternalServerError)
		return
	}

	var data requests.GetCreateOffer
	err = json.Unmarshal(req, &data)
	if err != nil {
		http.Error(w, "Problem for request", http.StatusInternalServerError)
		return
	}

	trip_id, err := crypt.Encrypt(string(req), key)
	if err != nil {
		http.Error(w, "Problem for request", 404)
		return
	}

	var cost int64 = int64(math.Sqrt((data.From.Lat-data.To.Lat)*(data.From.Lat-data.To.Lat)+
		(data.From.Lng-data.To.Lng)*(data.From.Lng-data.To.Lng)) + 200)

	offer, err := json.Marshal(responses.PullCreateOffer{Trip_id: trip_id,
		From:      responses.LatLngLiteral{Lat: data.From.Lat, Lng: data.From.Lng},
		To:        responses.LatLngLiteral{Lat: data.To.Lat, Lng: data.To.Lng},
		Client_id: data.Client_id, Price: responses.Money{Amount: cost, Currency: "RUB"}})

	w.WriteHeader(http.StatusOK)
	w.Write(offer)
}

func ParseOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := vars["trip_id"]
	if err != true {
		http.Error(w, "Couldn't parse trip_id", http.StatusBadRequest)
		return
	}

	str_offer, er := crypt.Decrypt(id, key)
	if er != nil {
		http.Error(w, "Trip don't find", 404)
		w.WriteHeader(404)
		return
	}

	offer := []byte(str_offer)

	w.WriteHeader(http.StatusOK)
	w.Write(offer)
}
