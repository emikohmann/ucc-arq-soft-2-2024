package users

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	hotelsDomain "search-api/domain/hotels"
	"strconv"
)

type Service interface {
	Search(ctx context.Context, query string, offset int, limit int) ([]hotelsDomain.Hotel, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) Search(c *gin.Context) {
	// Parse query from URL
	query := c.Query("q")

	// Parse offset from URL
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err),
		})
		return
	}

	// Parse limit from URL
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err),
		})
		return
	}

	// Invoke service
	hotels, err := controller.service.Search(c.Request.Context(), query, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error searching hotels: %s", err.Error()),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, hotels)
}
