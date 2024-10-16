package search

import (
	"context"
	"fmt"
	hotelsDAO "search-api/dao/hotels"
	hotelsDomain "search-api/domain/hotels"
)

type Repository interface {
	Search(ctx context.Context, query string, offset int, limit int) ([]hotelsDAO.Hotel, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service Service) Search(ctx context.Context, query string, offset int, limit int) ([]hotelsDomain.Hotel, error) {
	// Try to hotelsDAO in repository
	hotels, err := service.repository.Search(ctx, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error searching hotelsDAO: %s", err.Error())
	}

	// Convert
	result := make([]hotelsDomain.Hotel, 0)
	for _, hotel := range hotels {
		result = append(result, hotelsDomain.Hotel{
			ID:        hotel.ID,
			Name:      hotel.Name,
			Address:   hotel.Address,
			City:      hotel.City,
			State:     hotel.State,
			Rating:    hotel.Rating,
			Amenities: hotel.Amenities,
		})
	}

	// Send the result
	return result, nil
}
