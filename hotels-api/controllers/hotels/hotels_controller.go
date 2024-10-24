package hotels

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	hotelsDomain "hotels-api/domain/hotels"
	"net/http"
	"strings"
)

type Service interface {
	GetHotelByID(ctx context.Context, id string) (hotelsDomain.Hotel, error)
	Create(ctx context.Context, hotel hotelsDomain.Hotel) (string, error)
	Update(ctx context.Context, hotel hotelsDomain.Hotel) error
	Delete(ctx context.Context, id string) error
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) GetHotelByID(ctx *gin.Context) {
	// Validate ID param
	hotelID := strings.TrimSpace(ctx.Param("id"))

	// Get hotel by ID using the service
	hotel, err := controller.service.GetHotelByID(ctx.Request.Context(), hotelID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("error getting hotel: %s", err.Error()),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusOK, hotel)
}

func (controller Controller) Create(ctx *gin.Context) {
	// Parse hotel
	var hotel hotelsDomain.Hotel
	if err := ctx.ShouldBindJSON(&hotel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Create hotel
	id, err := controller.service.Create(ctx.Request.Context(), hotel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error creating hotel: %s", err.Error()),
		})
		return
	}

	// Send ID
	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (controller Controller) Update(ctx *gin.Context) {
	// Validate ID param
	id := strings.TrimSpace(ctx.Param("id"))

	// Parse hotel
	var hotel hotelsDomain.Hotel
	if err := ctx.ShouldBindJSON(&hotel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Set the ID from the URL to the hotel object
	hotel.ID = id

	// Update hotel
	if err := controller.service.Update(ctx.Request.Context(), hotel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error updating hotel: %s", err.Error()),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusOK, gin.H{
		"message": id,
	})
}

func (controller Controller) Delete(ctx *gin.Context) {
	// Validate ID param
	id := strings.TrimSpace(ctx.Param("id"))

	// Delete hotel
	if err := controller.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error deleting hotel: %s", err.Error()),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusOK, gin.H{
		"message": id,
	})
}
