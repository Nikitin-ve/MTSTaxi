package main

import (
	"client_service/internal/app/client_service/api/handlers"
	"client_service/internal/app/client_service/repository/taxi"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const shutdownTimeout = 60 * time.Second

const URI = "mongodb://127.0.0.1:27017"

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))

	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("Connected to MongoDB")

	defer func() {
		if client != nil {
			if err := client.Disconnect(context.Background()); err != nil {
				log.Println("Disconnecting error MongoDB:", err)
			}
		}
	}()

	serverMux := http.NewServeMux()

	my_taxi := handlers.MTSTaxiHandler{Taxi: &taxi.MTSTaxi{Client: client, Name: "trips"}}

	serverMux.HandleFunc("/trips", my_taxi.GPTrip)
	serverMux.HandleFunc("/trips/{trip_id}", my_taxi.GetTripID)
	serverMux.HandleFunc("/trip/{trip_id}/cancel", my_taxi.CancelTripID)

	server := http.Server{Addr: ":63343", Handler: serverMux}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	server.Shutdown(shutdownCtx)

}
