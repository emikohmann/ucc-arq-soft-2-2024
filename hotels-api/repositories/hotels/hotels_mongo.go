package hotels

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	hotelsDAO "hotels-api/dao/hotels"
	"log"
)

type MongoConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Database   string
	Collection string
}

type Mongo struct {
	client     *mongo.Client
	database   string
	collection string
}

const (
	connectionURI = "mongodb://%s:%s"
)

func NewMongo(config MongoConfig) Mongo {
	credentials := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	ctx := context.Background()
	uri := fmt.Sprintf(connectionURI, config.Host, config.Port)
	cfg := options.Client().ApplyURI(uri).SetAuth(credentials)

	client, err := mongo.Connect(ctx, cfg)
	if err != nil {
		log.Panicf("error connecting to mongo DB: %v", err)
	}

	return Mongo{
		client:     client,
		database:   config.Database,
		collection: config.Collection,
	}
}

func (repository Mongo) GetHotelByID(ctx context.Context, id string) (hotelsDAO.Hotel, error) {
	// Get from MongoDB
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return hotelsDAO.Hotel{}, fmt.Errorf("error converting id to mongo ID: %w", err)
	}
	result := repository.client.Database(repository.database).Collection(repository.collection).FindOne(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		return hotelsDAO.Hotel{}, fmt.Errorf("error finding document: %w", result.Err())
	}

	// Convert document to DAO
	var hotelDAO hotelsDAO.Hotel
	if err := result.Decode(&hotelDAO); err != nil {
		return hotelsDAO.Hotel{}, fmt.Errorf("error decoding result: %w", err)
	}
	return hotelDAO, nil
}

func (repository Mongo) Create(ctx context.Context, hotel hotelsDAO.Hotel) (string, error) {
	result, err := repository.client.Database(repository.database).Collection(repository.collection).InsertOne(ctx, hotel)
	if err != nil {
		return "", fmt.Errorf("error creating document: %w", err)
	}
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("error converting mongo ID to object ID")
	}
	return objectID.Hex(), nil
}
