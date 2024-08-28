package hotels

import (
	"hotels-api/dao/hotels"
)

type HotelsRepo interface {
	GetHotelByID(id int64) (hotels.HotelDAO, error)
}
