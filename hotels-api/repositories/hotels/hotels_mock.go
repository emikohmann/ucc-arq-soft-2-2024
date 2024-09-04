package hotels

import (
	"context"
	hotelsDAO "hotels-api/dao/hotels"
)

type Mock struct {
	docs map[int64]hotelsDAO.Hotel
}

func NewMock() Mock {
	return Mock{
		docs: map[int64]hotelsDAO.Hotel{
			1: {
				ID:      1,
				Name:    "Holiday Inn",
				Address: "Mock Address",
				City:    "Mock City",
				State:   "Mock State",
				Rating:  5,
				Amenities: []string{
					"Swimming Pool",
					"Free Wi-Fi",
				},
			},
		},
	}
}

func (repository Mock) GetHotelByID(ctx context.Context, id int64) (hotelsDAO.Hotel, error) {
	return repository.docs[id], nil
}
