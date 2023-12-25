package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"io"
	"log"
	"net/http"
	"trip_service/internal/app/trip_service/api/requests"
	"trip_service/internal/app/trip_service/repository/taxi"
)

type Repository struct {
	Rep *taxi.PostgreSQL
}

func NewReaderConfig(topic string) *kafka.ReaderConfig {
	return &kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
	}
}

func NewWriterKafka(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func Marshals[T requests.CreateTrip | requests.CancelTrip | requests.AcceptTrip | requests.StartEndTrip](m kafka.Message, trip T) (T, error) {
	req, err := json.Marshal(m)
	if err != nil {
		log.Printf("error marshal: %sn", err)
		return trip, err
	}

	err = json.Unmarshal(req, &trip)
	if err != nil {
		log.Printf("error unmarshal: %sn", err)
		return trip, err
	}
	return trip, nil
}

func (rep *Repository) Read1(r1 *kafka.Reader, ctx context.Context, w1 *kafka.Writer, w2 *kafka.Writer) {
	m, err := r1.FetchMessage(ctx)
	for err != nil {
		m, err = r1.FetchMessage(ctx)
	}

	if string(m.Key) == "Create" {
		var trip requests.CreateTrip
		trip, err := Marshals(m, trip)
		if err != nil {
			log.Printf("error while receiving message: %sn", err)
			return
		}
		req, err := http.NewRequest("GET", "http://localhost:63344/offers/"+trip.Data.Offer_id, nil)
		if err != nil {
			log.Printf("error while receiving message: %sn", err)
			return
		}
		req_byte, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("error while receiving message: %sn", err)
			return
		}
		var fromoffer requests.FromOffering
		err = json.Unmarshal(req_byte, &fromoffer)
		if err != nil {
			log.Printf("error while receiving message: %sn", err)
			return
		}
		req.Body.Close()
		full_ride := requests.ForDB{Trip_id: trip.Trip_id, Source: trip.Source,
			Type: trip.Type, Datacontenttype: trip.Datacontenttype,
			Time: trip.Time, Offer_id: trip.Data.Offer_id, Amount: fromoffer.Price.Amount,
			Currency: fromoffer.Price.Currency, Status: "CREATED", Lat_from: fromoffer.From.Lat,
			Lng_from: fromoffer.From.Lng, Lat_to: fromoffer.To.Lat, Lng_to: fromoffer.To.Lng}
		rep.Rep.AddNewTrip(full_ride)
		sm := requests.CreateTrip{Trip_id: trip.Trip_id, Source: trip.Source, Type: trip.Type,
			Datacontenttype: trip.Datacontenttype, Time: trip.Time, Data: requests.Offer1{Offer_id: trip.Data.Offer_id}}
		send_message, err := json.Marshal(sm)
		if err != nil {
			log.Printf("error marshal: %sn", err)
			return
		}
		err = w1.WriteMessages(ctx,
			kafka.Message{
				Value:   send_message,
				Headers: []protocol.Header{{Key: "client_write", Value: []byte(fmt.Sprint(m.Offset))}},
			},
		)
		if err != nil {
			log.Printf("error send message: %sn", err)
			return
		}
		err = w2.WriteMessages(ctx,
			kafka.Message{
				Value:   send_message,
				Headers: []protocol.Header{{Key: "driver_write", Value: []byte(fmt.Sprint(m.Offset))}},
			},
		)
		if err != nil {
			log.Printf("error send message: %sn", err)
			return
		}
	} else if string(m.Key) == "Cancel" {
		var trip requests.CancelTrip
		trip, err := Marshals(m, trip)
		if err != nil {
			log.Printf("error while receiving message: %sn", err)
			return
		}
		context.TODO()
	} else {
		log.Printf("error request")
	}
}

func (rep *Repository) Read2(r2 *kafka.Reader, ctx context.Context, w1 *kafka.Writer) {
	m, err := r2.FetchMessage(ctx)
	for err != nil {
		m, err = r2.FetchMessage(ctx)
	}

	if string(m.Key) == "Accept" {
		var trip requests.AcceptTrip
		trip, err := Marshals(m, trip)
		if err != nil {
			log.Printf("error while receiving message: %sn", err)
			return
		}
		sm := requests.StartEndTrip{Trip_id: trip.Trip_id, Source: trip.Source, Type: trip.Type,
			Datacontenttype: trip.Datacontenttype, Time: trip.Time, Data: requests.Offer4{Trip_id: trip.Data.Trip_id}}
		send_message, err := json.Marshal(sm)
		if err != nil {
			log.Printf("error marshal: %sn", err)
			return
		}
		err = w1.WriteMessages(ctx,
			kafka.Message{
				Value:   send_message,
				Headers: []protocol.Header{{Key: "client_write", Value: []byte(fmt.Sprint(m.Offset))}},
			},
		)
		if err != nil {
			log.Printf("error send message: %sn", err)
			return
		}
	} else if string(m.Key) == "Start" {
		var trip requests.StartEndTrip
		trip, err := Marshals(m, trip)
		if err != nil {
			log.Printf("error while receiving message: %sn", err)
			return
		}
		send_message, err := json.Marshal(trip)
		if err != nil {
			log.Printf("error marshal: %sn", err)
			return
		}
		err = w1.WriteMessages(ctx,
			kafka.Message{
				Value:   send_message,
				Headers: []protocol.Header{{Key: "client_write", Value: []byte(fmt.Sprint(m.Offset))}},
			},
		)
		if err != nil {
			log.Printf("error send message: %sn", err)
			return
		}
	} else if string(m.Key) == "End" {
		var trip requests.StartEndTrip
		trip, err := Marshals(m, trip)
		if err != nil {
			log.Printf("error while receiving message: %sn", err)
			return
		}
		send_message, err := json.Marshal(trip)
		if err != nil {
			log.Printf("error marshal: %sn", err)
			return
		}
		err = w1.WriteMessages(ctx,
			kafka.Message{
				Value:   send_message,
				Headers: []protocol.Header{{Key: "client_write", Value: []byte(fmt.Sprint(m.Offset))}},
			},
		)
		if err != nil {
			log.Printf("error send message: %sn", err)
			return
		}
	} else {
		log.Printf("error request")
	}
}

func (rep *Repository) StartKafka(ctx context.Context) {
	r1 := kafka.NewReader(*NewReaderConfig("client_read"))
	defer r1.Close()
	r2 := kafka.NewReader(*NewReaderConfig("driver_read"))
	defer r2.Close()

	w1 := NewWriterKafka("client_write")
	defer w1.Close()
	w2 := NewWriterKafka("driver_write")
	defer w2.Close()

	for {
		go rep.Read1(r1, ctx, w1, w2)
		go rep.Read2(r2, ctx, w1)
	}
}
