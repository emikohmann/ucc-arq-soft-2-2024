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
	repository.client.Set(key, hotel, repository.duration)
	return hotel.ID, nil
}

func (repository Cache) Update(ctx context.Context, hotel hotelsDAO.Hotel) error {
	key := fmt.Sprintf(keyFormat, hotel.ID)

	// Retrieve the current hotel data from the cache
	item := repository.client.Get(key)
	if item == nil {
		return fmt.Errorf("hotel with ID %s not found in cache", hotel.ID)
	}
	if item.Expired() {
		return fmt.Errorf("item with key %s is expired", key)
	}

	// Get the current hotel data
	currentHotel, ok := item.Value().(hotelsDAO.Hotel)
	if !ok {
		return fmt.Errorf("error converting item with key %s", key)
	}

	// Update only the fields that are non-zero or non-empty
	if hotel.Name != "" {
		currentHotel.Name = hotel.Name
	}
	if hotel.Address != "" {
		currentHotel.Address = hotel.Address
	}
	if hotel.City != "" {
		currentHotel.City = hotel.City
	}
	if hotel.State != "" {
		currentHotel.State = hotel.State
	}
	if hotel.Rating != 0 {
		currentHotel.Rating = hotel.Rating
	}
	if len(hotel.Amenities) > 0 {
		currentHotel.Amenities = hotel.Amenities
	}

	// Update the cache with the new hotel data and reset the expiration timer
	repository.client.Set(key, currentHotel, repository.duration)

	return nil
}

func (repository Cache) Delete(ctx context.Context, id string) error {
	key := fmt.Sprintf(keyFormat, id)
	// Remove the item from the cache
	repository.client.Delete(key)
	return nil
}
