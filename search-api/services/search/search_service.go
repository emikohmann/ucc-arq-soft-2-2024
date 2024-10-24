package search

import (
	"context"
	"fmt"
	hotelsDAO "search-api/dao/hotels"
	hotelsDomain "search-api/domain/hotels"
)

type Repository interface {
	Index(ctx context.Context, hotel hotelsDAO.Hotel) (string, error)
	Update(ctx context.Context, hotel hotelsDAO.Hotel) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int, offset int) ([]hotelsDAO.Hotel, error) // Updated signature
}

type ExternalRepository interface {
	GetHotelByID(ctx context.Context, id string) (hotelsDomain.Hotel, error)
}

type Service struct {
	repository Repository
	hotelsAPI  ExternalRepository
}

func NewService(repository Repository, hotelsAPI ExternalRepository) Service {
	return Service{
		repository: repository,
		hotelsAPI:  hotelsAPI,
	}
}

func (service Service) Search(ctx context.Context, query string, offset int, limit int) ([]hotelsDomain.Hotel, error) {
	// Call the repository's Search method
	hotelsDAOList, err := service.repository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error searching hotels: %w", err)
	}

	// Convert the dao layer hotels to domain layer hotels
	hotelsDomainList := make([]hotelsDomain.Hotel, 0)
	for _, hotel := range hotelsDAOList {
		hotelsDomainList = append(hotelsDomainList, hotelsDomain.Hotel{
			ID:        hotel.ID,
			Name:      hotel.Name,
			Address:   hotel.Address,
			City:      hotel.City,
			State:     hotel.State,
			Rating:    hotel.Rating,
			Amenities: hotel.Amenities,
		})
	}

	return hotelsDomainList, nil
}

func (service Service) HandleHotelNew(hotelNew hotelsDomain.HotelNew) {
	switch hotelNew.Operation {
	case "CREATE", "UPDATE":
		// Fetch hotel details from the local service
		hotel, err := service.hotelsAPI.GetHotelByID(context.Background(), hotelNew.HotelID)
		if err != nil {
			fmt.Printf("Error getting hotel (%s) from API: %v\n", hotelNew.HotelID, err)
			return
		}

		hotelDAO := hotelsDAO.Hotel{
			ID:        hotel.ID,
			Name:      hotel.Name,
			Address:   hotel.Address,
			City:      hotel.City,
			State:     hotel.State,
			Rating:    hotel.Rating,
			Amenities: hotel.Amenities,
		}

		// Handle Index operation
		if hotelNew.Operation == "CREATE" {
			if _, err := service.repository.Index(context.Background(), hotelDAO); err != nil {
				fmt.Printf("Error indexing hotel (%s): %v\n", hotelNew.HotelID, err)
			} else {
				fmt.Println("Hotel indexed successfully:", hotelNew.HotelID)
			}
		} else { // Handle Update operation
			if err := service.repository.Update(context.Background(), hotelDAO); err != nil {
				fmt.Printf("Error updating hotel (%s): %v\n", hotelNew.HotelID, err)
			} else {
				fmt.Println("Hotel updated successfully:", hotelNew.HotelID)
			}
		}

	case "DELETE":
		// Call Delete method directly since no hotel details are needed
		if err := service.repository.Delete(context.Background(), hotelNew.HotelID); err != nil {
			fmt.Printf("Error deleting hotel (%s): %v\n", hotelNew.HotelID, err)
		} else {
			fmt.Println("Hotel deleted successfully:", hotelNew.HotelID)
		}

	default:
		fmt.Printf("Unknown operation: %s\n", hotelNew.Operation)
	}
}
