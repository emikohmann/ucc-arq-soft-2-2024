package repositories

import "backend/dao"

type HotelsMock struct{}

func NewHotelsMock() HotelsMock {
	return HotelsMock{}
}
func (HotelsMock) GetHotelByID(id int64) (dao.HotelDAO, error) {

	return dao.HotelDAO{
		ID:        id,
		Name:      "HotelMock",
		Address:   "Mock Address",
		City:      "Mock City",
		State:     "Mock Stats",
		Rating:    5,
		Amenities: nil,
	}, nil
}
