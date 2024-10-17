package hotels

import (
	"context"
	"fmt"
	hotelsDAO "hotels-api/dao/hotels"
	hotelsDomain "hotels-api/domain/hotels"
)

type Repository interface {
	GetHotelByID(ctx context.Context, id string) (hotelsDAO.Hotel, error)
	Create(ctx context.Context, hotel hotelsDAO.Hotel) (string, error)
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
		fmt.Println(fmt.Sprintf("Hotel %s not found in secondary repository", id))
		// Get hotel from main repository
		hotelDAO, err = service.mainRepository.GetHotelByID(ctx, id)
		if err != nil {
			fmt.Println(fmt.Sprintf("Hotel %s not found in main repository", id))
			return hotelsDomain.Hotel{}, fmt.Errorf("error getting hotel from repository: %v", err)
		}
		// Set ID from main repository to use in the rest of the repositories
		fmt.Println(fmt.Sprintf("Hotel %s retrieved from main repository", id))
		if _, err := service.cacheRepository.Create(ctx, hotelDAO); err != nil {
			fmt.Println(fmt.Sprintf("Warning: error creating hotel in cache: %w", err))
		} else {
			fmt.Println(fmt.Sprintf("Hotel %s saved in secondary repository", id))
		}
	} else {
		fmt.Println(fmt.Sprintf("Hotel %s retrieved from secondary repository", id))
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

func (service Service) Create(ctx context.Context, hotel hotelsDomain.Hotel) (string, error) {
	record := hotelsDAO.Hotel{
		Name:      hotel.Name,
		Address:   hotel.Address,
		City:      hotel.City,
		State:     hotel.State,
		Rating:    hotel.Rating,
		Amenities: hotel.Amenities,
	}
	id, err := service.mainRepository.Create(ctx, record)
	if err != nil {
		return "", fmt.Errorf("error creating hotel in main repository: %w", err)
	}
	// Set ID from main repository to use in the rest of the repositories
	record.ID = id
	if _, err := service.cacheRepository.Create(ctx, record); err != nil {
		fmt.Println(fmt.Sprintf("error creating hotel in cache: %w", err))
	} else {
		fmt.Println("saved in secondary repository")
	}
	if err := service.eventsQueue.Publish(hotelsDomain.HotelNew{
		Operation: "GET",
		HotelID:   id,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Error publishing hotel new: %w", err))
	}
	fmt.Println(fmt.Sprintf("Published to rabbitMQ hotelID %s Operation GET", id))
	return id, err
}
