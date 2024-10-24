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
	Search(ctx context.Context, query string) ([]hotelsDAO.Hotel, error)
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
	return nil, nil
}

func (service Service) HandleHotelNew(hotelNew hotelsDomain.HotelNew) {
	fmt.Println("Received new", hotelNew)
}
