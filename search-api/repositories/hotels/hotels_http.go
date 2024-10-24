package hotels

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	hotelsDomain "search-api/domain/hotels"
)

type HTTPConfig struct {
	Host string
	Port string
}

type HTTP struct {
	baseURL func(hotelID string) string
}

func NewHTTP(config HTTPConfig) HTTP {
	return HTTP{
		baseURL: func(hotelID string) string {
			return fmt.Sprintf("http://%s:%s/hotels/%s", config.Host, config.Port, hotelID)
		},
	}
}

func (repository HTTP) GetHotelByID(ctx context.Context, id string) (hotelsDomain.Hotel, error) {
	resp, err := http.Get(repository.baseURL(id))
	if err != nil {
		return hotelsDomain.Hotel{}, fmt.Errorf("Error fetching hotel (%s): %w\n", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return hotelsDomain.Hotel{}, fmt.Errorf("Failed to fetch hotel (%s): received status code %d\n", id, resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return hotelsDomain.Hotel{}, fmt.Errorf("Error reading response body for hotel (%s): %w\n", id, err)
	}

	// Unmarshal the hotel details into the hotel struct
	var hotel hotelsDomain.Hotel
	if err := json.Unmarshal(body, &hotel); err != nil {
		return hotelsDomain.Hotel{}, fmt.Errorf("Error unmarshaling hotel data (%s): %w\n", id, err)
	}

	return hotel, nil
}
