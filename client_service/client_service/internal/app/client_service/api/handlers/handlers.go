package handlers

import (
	"client_service/internal/app/client_service/api/requests"
	"client_service/internal/app/client_service/api/responses"
	"client_service/internal/app/client_service/repository/taxi"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"io"
	"log"
	"net/http"
)

type TaxiHandler interface {
	GPTrip(w http.ResponseWriter, r *http.Request)
	GetTripID(w http.ResponseWriter, r *http.Request)
	CancelTripID(w http.ResponseWriter, r *http.Request)
}

type MTSTaxiHandler struct {
	Taxi *taxi.MTSTaxi
}

func (b *MTSTaxiHandler) GPTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// получаю user_id и offer_id, отправляюсь в OfferingSvc, если некорректный offer_id, возвращаю ошибку, иначе иду
		// в TripSvc и получаю одобрение, передаю клиенту

		req, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Problem for request", 400)
			w.WriteHeader(400)
			return
		}
		var off responses.FromOffer
		err = json.Unmarshal(req, &off)
		if err != nil {
			http.Error(w, "Problem for unmarshal", 400)
			w.WriteHeader(400)
			return
		}
		ans, err := http.NewRequest("GET", "http://localhost:63344/offers/"+off.Offer_id, nil)
		if err != nil {
			http.Error(w, "Incorrect offer id", 400)
			w.WriteHeader(400)
			return
		}
		ans_byte, err := io.ReadAll(ans.Body)
		if err != nil {
			http.Error(w, "Incorrect offer id", 400)
			w.WriteHeader(400)
			return
		}
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "client_write",
		})
		wr := kafka.Writer(kafka.Writer{
			Addr:     kafka.TCP("localhost:9092"),
			Topic:    "client_read",
			Balancer: &kafka.LeastBytes{},
		})
		err = wr.WriteMessages(context.Background(),
			kafka.Message{
				Value:   ans_byte,
				Headers: []protocol.Header{{Key: "Create", Value: []byte(fmt.Sprint(0))}},
			},
		)
		m, err := r.FetchMessage(context.Background())
		for err != nil {
			m, err = r.FetchMessage(context.Background())
		}
		reqest, err := json.Marshal(m)
		if err != nil {
			http.Error(w, "Incorrect offer id", 400)
			w.WriteHeader(400)
			return
		}

		var trip requests.CreatedTrip

		err = json.Unmarshal(reqest, &trip)
		if err != nil {
			log.Printf("error unmarshal: %sn", 400)
			w.WriteHeader(400)
			return
		}
		db_abs := requests.ForDB{Trip_id: trip.Trip_id, Source: trip.Source, Type: trip.Type,
			Datacontenttype: trip.Datacontenttype, Time: trip.Time, Offer_id: trip.Data.Offer_id,
			Amount: trip.Data.Price.Amount, Currency: trip.Data.Price.Currency, Status: trip.Data.Status,
			Lat_from: trip.Data.From.Lat, Lng_from: trip.Data.From.Lng, Lng_to: trip.Data.To.Lng,
			Lat_to: trip.Data.To.Lat}
		err = b.Taxi.CreateTrip(db_abs)
		if err != nil {
			log.Printf("error BD: %sn", 400)
			w.WriteHeader(400)
			return
		}
	}
	// получаю запрос от клиента дать ему список всех его поездок, выдаю из бд

	req, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Problem for request", http.StatusInternalServerError)
		return
	}

	var user_id requests.GetUserID
	err = json.Unmarshal(req, &user_id)
	if err != nil {
		http.Error(w, "I can't parse a byte slice", http.StatusInternalServerError)
		return
	}

	ans, err := b.Taxi.GetAllTrip(user_id.User_id)
	if err != nil {
		http.Error(w, "An error occurred with the search for your trips to the storage", http.StatusInternalServerError)
		return
	}

	answer, err := json.Marshal(ans)
	if err != nil {
		http.Error(w, "An error occurred while gluing", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func (b *MTSTaxiHandler) GetTripID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// отправляю поездку по ее ID

		vars := mux.Vars(r)
		id, err := vars["trip_id"]
		if err != true {
			http.Error(w, "Couldn't parse trip_id", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		a, er := b.Taxi.GetTrip(id)
		if er != nil {
			http.Error(w, "I can't find trip_id", http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ans, er := json.Marshal(a)
		if er != nil {
			http.Error(w, "An error occurred while gluing", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(ans)
	}
	http.Error(w, "Wrong request", 500)
}

func (b *MTSTaxiHandler) CancelTripID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// удаляю поездку из БД по ее ID

		vars := mux.Vars(r)
		id, err := vars["trip_id"]
		if err != true {
			http.Error(w, "Couldn't parse trip_id", http.StatusBadRequest)
			return
		}

		er := b.Taxi.DeleteResult(id)
		if er != nil {
			http.Error(w, "I can't find trip_id", http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
	http.Error(w, "Wrong request", 500)
}
