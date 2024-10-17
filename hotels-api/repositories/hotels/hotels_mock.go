package hotels

import (
	"context"
	"github.com/google/uuid"
	hotelsDAO "hotels-api/dao/hotels"
)

type Mock struct {
	docs map[string]hotelsDAO.Hotel
}

func NewMock() Mock {
	return Mock{
		docs: make(map[string]hotelsDAO.Hotel),
	}
}

func (repository Mock) GetHotelByID(ctx context.Context, id string) (hotelsDAO.Hotel, error) {
	return repository.docs[id], nil
}

func (repository Mock) Create(ctx context.Context, hotel hotelsDAO.Hotel) (string, error) {
	id := uuid.New().String()
	hotel.ID = uuid.New().String()
	repository.docs[id] = hotel
	return id, nil
}
