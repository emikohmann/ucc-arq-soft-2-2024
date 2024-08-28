package repositories

import (
	"backend/dao"
)

type HotelsRepo interface {
	GetHotelByID(id int64) (dao.HotelDAO, error)
}
