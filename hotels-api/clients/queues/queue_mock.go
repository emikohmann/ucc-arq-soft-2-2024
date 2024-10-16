package queues

import "hotels-api/domain/hotels"

type Mock struct{}

func NewMock() Mock {
	return Mock{}
}

func (Mock) Publish(hotelNew hotels.HotelNew) error {
	return nil
}
