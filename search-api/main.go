package main

import (
	"github.com/gin-gonic/gin"
	"log"
	controllers "search-api/controllers/search"
	search "search-api/repositories/hotels"
	services "search-api/services/search"
)

func main() {
	// Solr
	solrConfig := search.SolrConfig{
		BaseURL:    "",
		Collection: "hotels",
	}
	solrRepo := search.NewSolr(solrConfig)

	// Services
	service := services.NewService(solrRepo)

	// Handlers
	controller := controllers.NewController(service)

	// Create router
	router := gin.Default()

	// URL mappings
	// /hotels/search?q=sheraton&limit=10&offset=0
	router.GET("/hotels/search", controller.Search)

	// Run application
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running application: %v", err)
	}
}
