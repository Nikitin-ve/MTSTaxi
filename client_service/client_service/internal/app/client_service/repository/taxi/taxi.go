package taxi

import (
	"client_service/internal/app/client_service/api/requests"
	"client_service/pkg/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MTSTaxi struct {
	Client *mongo.Client
	Name   string
}

func (a *MTSTaxi) GetAllTrip(user_id string) ([]*models.Trip, error) {
	filter := bson.D{{"User_id", user_id}}
	var results []*models.Trip

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := a.Client.Database(a.Name).Collection("trips")

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem models.Trip
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)

	return results, nil
}

func (a *MTSTaxi) GetTrip(trip_id string) (models.Trip, error) {
	filter := bson.D{{"Trip_id", trip_id}}
	var result models.Trip

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := a.Client.Database(a.Name).Collection("trips")

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return models.Trip{}, err
	}

	return result, nil
}

func (a *MTSTaxi) DeleteResult(trip_id string) error {
	filter := bson.D{{"Trip_id", trip_id}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := a.Client.Database(a.Name).Collection("trips")

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (a *MTSTaxi) CreateTrip(db_abs requests.ForDB) error {
	collection := a.Client.Database(a.Name).Collection("trips")

	_, err := collection.InsertOne(context.Background(), db_abs)
	if err != nil {
		return err
	}
	return nil
}
