package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"search-api/clients/queues"
	controllers "search-api/controllers/search"
	repositories "search-api/repositories/hotels"
	services "search-api/services/search"
)

func main() {
	// Solr
	solrRepo := repositories.NewSolr(repositories.SolrConfig{
		Host:       "solr",   // Solr host
		Port:       "8983",   // Solr port
		Collection: "hotels", // Collection name
	})

	// Rabbit
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "root",
		Password:  "root",
		QueueName: "hotels-news",
	})

	// Hotels API
	hotelsAPI := repositories.NewHTTP(repositories.HTTPConfig{
		Host: "hotels-api",
		Port: "8081",
	})

	// Services
	service := services.NewService(solrRepo, hotelsAPI)

	// Controllers
	controller := controllers.NewController(service)

	// Launch rabbit consumer
	if err := eventsQueue.StartConsumer(service.HandleHotelNew); err != nil {
		log.Fatalf("Error running consumer: %v", err)
	}

	// Create router
	router := gin.Default()
	router.GET("/search", controller.Search)
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
