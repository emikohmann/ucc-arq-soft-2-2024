package hotels

import (
	"context"
	"fmt"
	"github.com/karlseguin/ccache"
	_ "github.com/karlseguin/ccache"
	hotelsDAO "hotels-api/dao/hotels"
)

const (
	keyFormat = "hotel:%d"
)

type CacheConfig struct {
	MaxSize      int64
	ItemsToPrune uint32
}

type Cache struct {
	client *ccache.Cache
}

func NewCache(config CacheConfig) Cache {
	client := ccache.New(ccache.Configure().
		MaxSize(config.MaxSize).
		ItemsToPrune(config.ItemsToPrune))
	return Cache{
		client: client,
	}
}

func (repo Cache) GetHotelByID(
	ctx context.Context, id int64) (
	hotelsDAO.Hotel, error) {

	key := fmt.Sprintf(keyFormat, id)
	item := repo.client.Get(key)
	hotelDAO, ok := item.Value().(hotelsDAO.Hotel)
	if !ok {
		return hotelsDAO.Hotel{},
			fmt.Errorf("error converting item with key %s", key)
	}
	return hotelDAO, nil
}
