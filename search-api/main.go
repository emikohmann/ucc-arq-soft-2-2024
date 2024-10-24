package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"search-api/clients/queues"
	controllers "search-api/controllers/search"
	search "search-api/repositories/hotels"
	services "search-api/services/search"
)

func main() {
	// Solr
	solrConfig := search.SolrConfig{
		BaseURL:    "localhost",
		Collection: "hotels",
	}
	solrRepo := search.NewSolr(solrConfig)

	// Rabbit
	rabbitConfig := queues.RabbitConfig{
		Username:  "root",
		Password:  "root",
		Host:      "localhost",
		Port:      "5672",
		QueueName: "hotels-news",
	}

	// Dependencies
	eventsQueue := queues.NewRabbit(rabbitConfig)
	service := services.NewService(solrRepo)
	controller := controllers.NewController(service)

	// Launch rabbit consumer
	if err := eventsQueue.StartConsumer(service.HandleHotelNew); err != nil {
		log.Fatalf("Error running consumer: %v", err)
	}

	// Create router
	router := gin.Default()
	router.GET("/hotels/search", controller.Search)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
