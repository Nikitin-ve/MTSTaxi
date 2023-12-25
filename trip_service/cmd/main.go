package main

import (
	"context"
	"log"
	"os"
	"trip_service/internal/app/trip_service/repository/taxi"
	"trip_service/pkg/kafka"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	urlExample := "postgres://gerk:1881@localhost:5432/trip"
	db, err := pgxpool.New(context.Background(), urlExample)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx := context.Background()

	DB := kafka.Repository{Rep: &taxi.PostgreSQL{DB: db}}
	DB.StartKafka(ctx)
}
