package main

import (
	"github.com/gin-gonic/gin"
	"hotels-api/clients/queues"
	controllers "hotels-api/controllers/hotels"
	repositories "hotels-api/repositories/hotels"
	services "hotels-api/services/hotels"
	"log"
	"time"
)

func main() {
	// Local cache
	cacheConfig := repositories.CacheConfig{
		MaxSize:      100000,
		ItemsToPrune: 100,
		Duration:     30 * time.Second,
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
		Username:  "user",
		Password:  "password",
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
	router.POST("/hotels", controller.Create)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("error running application: %w", err)
	}
}
