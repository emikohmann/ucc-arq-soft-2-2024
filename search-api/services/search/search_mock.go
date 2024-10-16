package search

import (
	"context"
	hotelsDomain "search-api/domain/hotels"
)

type Mock struct{}

func NewMock() Mock {
	return Mock{}
}

func (service Mock) Search(ctx context.Context, query string, offset int, limit int) ([]hotelsDomain.Hotel, error) {
	//TODO implement me
	panic("implement me")
}
