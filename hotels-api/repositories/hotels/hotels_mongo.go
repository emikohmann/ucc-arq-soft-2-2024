package hotels

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotels-api/dao/hotels"
	"time"
)

type HotelsMongo struct {
	Client *mongo.Client
}

func NewHotelsMongo() HotelsMongo {
	// Creamos un contexto
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second)
	defer cancel()

	// Creamos las opciones
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{
			Username: "root",
			Password: "root",
		})

	// Conectamos el cliente
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error: ", err)
		return HotelsMongo{}
	}

	return HotelsMongo{
		Client: client,
	}
}

func (repo HotelsMongo) GetHotelByID(id int64) (
	hotels.HotelDAO, error) {

	ctx := context.Background()

	result := repo.Client.
		Database("hotels").
		Collection("hotels").
		FindOne(ctx, bson.M{"id": id})
	if result.Err() != nil {
		fmt.Println("Error: ", result.Err())
		return hotels.HotelDAO{}, result.Err()
	}

	var hotelDAO hotels.HotelDAO
	err := result.Decode(&hotelDAO)
	if err != nil {
		fmt.Println("Error: ", result.Err())
		return hotels.HotelDAO{}, result.Err()
	}

	return hotelDAO, nil
}