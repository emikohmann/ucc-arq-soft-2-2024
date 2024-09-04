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
	mongoConfig := hotelsRepository.MongoConfig{
		Host:       "localhost",
		Port:       "27017",
		Username:   "root",
		Password:   "root",
		Database:   "hotels-api",
		Collection: "hotels",
	}

	// Dependencies
	repository := hotelsRepository.NewMongo(mongoConfig)
	service := hotelsService.NewService(repository)
	controller := hotelsController.NewController(service)

	// Router
	router := gin.Default()
	router.GET("/hotels/:id", controller.GetHotelByID)
}
