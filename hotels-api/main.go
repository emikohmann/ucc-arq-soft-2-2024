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
	cacheRepository := repositories.NewCache(repositories.CacheConfig{
		MaxSize:      100000,
		ItemsToPrune: 100,
		Duration:     30 * time.Second,
	})

	// Mongo
	mainRepository := repositories.NewMongo(repositories.MongoConfig{
		Host:       "mongo",
		Port:       "27017",
		Username:   "root",
		Password:   "root",
		Database:   "hotels-api",
		Collection: "hotels",
	})

	// Rabbit
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "root",
		Password:  "root",
		QueueName: "hotels-news",
	})

	// Services
	service := services.NewService(mainRepository, cacheRepository, eventsQueue)

	// Controllers
	controller := controllers.NewController(service)

	// Router
	router := gin.Default()
	router.GET("/hotels/:id", controller.GetHotelByID)
	router.POST("/hotels", controller.Create)
	router.PUT("/hotels/:id", controller.Update)
	router.DELETE("/hotels/:id", controller.Delete)
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("error running application: %w", err)
	}
}
