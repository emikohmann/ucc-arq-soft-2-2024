package hotels

import (
	"context"
	"fmt"
	"hotels-api/dao/hotels"
	hotelsDomain "hotels-api/domain/hotels"
)

type Repository interface {
	GetHotelByID(ctx context.Context, id string) (hotels.Hotel, error)
}

type Queue interface {
	Publish(hotelNew hotelsDomain.HotelNew) error
}

type Service struct {
	mainRepository  Repository
	cacheRepository Repository
	eventsQueue     Queue
}

func NewService(mainRepository Repository, cacheRepository Repository, eventsQueue Queue) Service {
	return Service{
		mainRepository:  mainRepository,
		cacheRepository: cacheRepository,
		eventsQueue:     eventsQueue,
	}
}

func (service Service) GetHotelByID(ctx context.Context, id string) (hotelsDomain.Hotel, error) {
	hotelDAO, err := service.cacheRepository.GetHotelByID(ctx, id)
	if err != nil {
		// Get hotel from main repository
		hotelDAO, err = service.mainRepository.GetHotelByID(ctx, id)
		if err != nil {
			return hotelsDomain.Hotel{}, fmt.Errorf("error getting hotel from repository: %v", err)
		}

		// TODO: service.cacheRepository.CreateHotel
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

/**

// Example of notifying to the events queue
go func() {
	if err := service.eventsQueue.Publish(hotelsDomain.HotelNew{
		Operation: "GET",
		HotelID:   hotelDAO.ID,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Error publishing hotel new: %w", err))
	}
}()

*/
