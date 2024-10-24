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
	Update(ctx context.Context, hotel hotelsDAO.Hotel) error
	Delete(ctx context.Context, id string) error
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
		// Set ID from main repository to use in the rest of the repositories
		if _, err := service.cacheRepository.Create(ctx, hotelDAO); err != nil {
			return hotelsDomain.Hotel{}, fmt.Errorf("error creating hotel in cache: %w", err)
		}
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
		return "", fmt.Errorf("error creating hotel in cache: %w", err)
	}
	if err := service.eventsQueue.Publish(hotelsDomain.HotelNew{
		Operation: "CREATE",
		HotelID:   id,
	}); err != nil {
		return "", fmt.Errorf("error publishing hotel new: %w", err)
	}

	return id, nil
}

func (service Service) Update(ctx context.Context, hotel hotelsDomain.Hotel) error {
	// Convert domain model to DAO model
	record := hotelsDAO.Hotel{
		ID:        hotel.ID,
		Name:      hotel.Name,
		Address:   hotel.Address,
		City:      hotel.City,
		State:     hotel.State,
		Rating:    hotel.Rating,
		Amenities: hotel.Amenities,
	}

	// Update the hotel in the main repository
	err := service.mainRepository.Update(ctx, record)
	if err != nil {
		return fmt.Errorf("error updating hotel in main repository: %w", err)
	}

	// Try to update the hotel in the cache repository
	if err := service.cacheRepository.Update(ctx, record); err != nil {
		return fmt.Errorf("error updating hotel in cache: %w", err)
	}

	// Publish an event for the update operation
	if err := service.eventsQueue.Publish(hotelsDomain.HotelNew{
		Operation: "UPDATE",
		HotelID:   hotel.ID,
	}); err != nil {
		return fmt.Errorf("error publishing hotel update: %w", err)
	}

	return nil
}

func (service Service) Delete(ctx context.Context, id string) error {
	// Delete the hotel from the main repository
	err := service.mainRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting hotel from main repository: %w", err)
	}

	// Try to delete the hotel from the cache repository
	if err := service.cacheRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting hotel from cache: %w", err)
	}

	// Publish an event for the delete operation
	if err := service.eventsQueue.Publish(hotelsDomain.HotelNew{
		Operation: "DELETE",
		HotelID:   id,
	}); err != nil {
		return fmt.Errorf("error publishing hotel delete: %w", err)
	}

	return nil
}
