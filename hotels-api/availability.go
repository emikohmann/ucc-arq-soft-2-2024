package main

import "sync"

type Availability struct {
	HotelID   string
	Available bool
}

func GetAvailability(hotelIDs []string) map[string]bool {
	result := make(map[string]bool)
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(hotelIDs))
	ch := make(chan Availability)
	go func() {
		for {
			availability := <-ch
			result[availability.HotelID] = availability.Available
		}
	}()
	for _, hotelID := range hotelIDs {
		go IsAvailableAsync(hotelID, &waitGroup, ch)
	}
	waitGroup.Wait()
	return result
}

func IsAvailableAsync(hotelID string, group *sync.WaitGroup, ch chan Availability) {
	defer group.Done()
	ch <- Availability{
		HotelID:   hotelID,
		Available: IsAvailable(hotelID),
	}
}

func IsAvailable(hotelID string) bool {
	// implement
	return true
}
