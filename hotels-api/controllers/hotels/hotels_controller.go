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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting hotel: %s", err.Error()),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusOK, hotel)
}
