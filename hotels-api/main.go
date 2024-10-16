package main

import (
	"github.com/gin-gonic/gin"
	"hotels-api/clients/queues"
	controllers "hotels-api/controllers/hotels"
	repositories "hotels-api/repositories/hotels"
	services "hotels-api/services/hotels"
)

func main() {
	// Local cache
	cacheConfig := repositories.CacheConfig{
		MaxSize:      100000,
		ItemsToPrune: 100,
	}

	// Mongo
	mongoConfig := repositories.MongoConfig{
		Host:       "localhost",
		Port:       "27017",
		Username:   "root",
		Password:   "root",
		Database:   "hotels-api",
		Collection: "hotels",
	}

	// Rabbit
	rabbitConfig := queues.RabbitConfig{
		Username:  "guest",
		Password:  "guest",
		Host:      "localhost",
		Port:      "5672",
		QueueName: "hotels-news",
	}
	eventsQueue := queues.NewRabbit(rabbitConfig)

	// Dependencies
	mainRepository := repositories.NewMongo(mongoConfig)
	cacheRepository := repositories.NewCache(cacheConfig)
	service := services.NewService(mainRepository, cacheRepository, eventsQueue)
	controller := controllers.NewController(service)

	// Router
	router := gin.Default()
	router.GET("/hotels/:id", controller.GetHotelByID)
}
