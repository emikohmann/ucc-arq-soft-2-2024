package hotels

import (
	"context"
	"fmt"
	"github.com/karlseguin/ccache"
	hotelsDAO "hotels-api/dao/hotels"
	"time"
)

const (
	keyFormat = "hotel:%s"
)

type CacheConfig struct {
	MaxSize      int64
	ItemsToPrune uint32
	Duration     time.Duration
}

type Cache struct {
	client   *ccache.Cache
	duration time.Duration
}

func NewCache(config CacheConfig) Cache {
	client := ccache.New(ccache.Configure().
		MaxSize(config.MaxSize).
		ItemsToPrune(config.ItemsToPrune))
	return Cache{
		client:   client,
		duration: config.Duration,
	}
}

func (repository Cache) GetHotelByID(ctx context.Context, id string) (hotelsDAO.Hotel, error) {
	key := fmt.Sprintf(keyFormat, id)
	item := repository.client.Get(key)
	fmt.Println(key)
	if item == nil {
		return hotelsDAO.Hotel{}, fmt.Errorf("not found item with key %s", key)
	}
	if item.Expired() {
		return hotelsDAO.Hotel{}, fmt.Errorf("item with key %s is expired", key)
	}
	hotelDAO, ok := item.Value().(hotelsDAO.Hotel)
	if !ok {
		return hotelsDAO.Hotel{}, fmt.Errorf("error converting item with key %s", key)
	}
	return hotelDAO, nil
}

func (repository Cache) Create(ctx context.Context, hotel hotelsDAO.Hotel) (string, error) {
	key := fmt.Sprintf(keyFormat, hotel.ID)
	fmt.Println("saving with duration", repository.duration)
	repository.client.Set(key, hotel, repository.duration)
	return hotel.ID, nil
}
