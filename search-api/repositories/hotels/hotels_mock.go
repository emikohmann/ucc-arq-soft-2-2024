package hotels

import (
	"search-api/dao/hotels"
)

type Mock struct {
	data map[int64]hotels.Hotel
}

func NewMock() Mock {
	return Mock{
		data: make(map[int64]hotels.Hotel),
	}
}
