package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"offering_service/internal/app/offering_service/api/handlers"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = 60 * time.Second

func main() {
	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/offers", handlers.CreateOffer)
	serverMux.HandleFunc("/offers/{offer_id}", handlers.ParseOffer)

	server := http.Server{Addr: ":63344", Handler: serverMux}

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
