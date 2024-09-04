package hotels

import (
	"context"
	"fmt"
	"hotels-api/dao/hotels"
	hotelsDomain "hotels-api/domain/hotels"
)

type Repository interface {
	GetHotelByID(ctx context.Context, id int64) (hotels.Hotel, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service Service) GetHotelByID(ctx context.Context, id int64) (hotelsDomain.Hotel, error) {
	// Get hotel from repository
	hotelDAO, err := service.repository.GetHotelByID(ctx, id)
	if err != nil {
		return hotelsDomain.Hotel{}, fmt.Errorf("error getting hotel from repository: %v", err)
	}

	// Convert DAO to DTO
	return hotelsDomain.Hotel{
		ID:        hotelDAO.ID,
		Name:      hotelDAO.Name,
		Address:   hotelDAO.Address,
		City:      hotelDAO.City,
		State:     hotelDAO.State,
		Rating:    hotelDAO.Rating,
		Amenities: hotelDAO.Amenities,
	}, nil
}
