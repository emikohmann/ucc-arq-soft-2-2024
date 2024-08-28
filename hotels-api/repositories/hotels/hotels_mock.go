package hotels

import (
	"hotels-api/dao/hotels"
)

type HotelsMock struct{}

func NewHotelsMock() HotelsMock {
	return HotelsMock{}
}
func (HotelsMock) GetHotelByID(id int64) (hotels.HotelDAO, error) {

	return hotels.HotelDAO{
		ID:        id,
		Name:      "HotelMock",
		Address:   "Mock Address",
		City:      "Mock City",
		State:     "Mock Stats",
		Rating:    5,
		Amenities: nil,
	}, nil
}
