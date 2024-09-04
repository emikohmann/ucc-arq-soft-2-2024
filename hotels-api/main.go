package main

import (
	"github.com/gin-gonic/gin"
	hotelsController "hotels-api/controllers/hotels"
	hotelsRepository "hotels-api/repositories/hotels"
	hotelsService "hotels-api/services/hotels"
)

type Controller interface {
	GetHotelByID(ctx *gin.Context)
}

func main() {
	// Config
	cacheConfig := hotelsRepository.CacheConfig{
		MaxSize:      100000,
		ItemsToPrune: 100,
	}

	mongoConfig := hotelsRepository.MongoConfig{
		Host:       "localhost",
		Port:       "27017",
		Username:   "root",
		Password:   "root",
		Database:   "hotels-api",
		Collection: "hotels",
	}

	// Dependencies
	mainRepository := hotelsRepository.NewMongo(mongoConfig)
	cacheRepository := hotelsRepository.NewCache(cacheConfig)
	service := hotelsService.NewService(mainRepository, cacheRepository)
	controller := hotelsController.NewController(service)

	// Router
	router := gin.Default()
	router.GET("/hotels/:id", controller.GetHotelByID)
}
