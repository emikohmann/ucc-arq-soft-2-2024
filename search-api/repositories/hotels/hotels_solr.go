package hotels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stevenferrer/solr-go"
	"search-api/dao/hotels"
)

type SolrConfig struct {
	BaseURL    string
	Collection string
}

type Solr struct {
	Client     *solr.JSONClient
	Collection string
}

// NewSolr initializes a new Solr client
func NewSolr(config SolrConfig) (*Solr, error) {
	client := solr.NewJSONClient(config.BaseURL)

	return &Solr{
		Client:     client,
		Collection: config.Collection,
	}, nil
}

// Index adds a new hotel document to the Solr collection
func (searchEngine Solr) Index(ctx context.Context, hotel hotels.Hotel) (string, error) {
	// Prepare the document for Solr
	doc := map[string]interface{}{
		"id":        hotel.ID,
		"name":      hotel.Name,
		"address":   hotel.Address,
		"city":      hotel.City,
		"state":     hotel.State,
		"rating":    hotel.Rating,
		"amenities": hotel.Amenities,
	}

	// Index the document in Solr
	body, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("error marshaling hotel document: %w", err)
	}

	// Index the document in Solr
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error indexing hotel: %w", err)
	}
	if resp.Error != nil {
		return "", fmt.Errorf("failed to index hotel: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return "", fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return hotel.ID, nil
}

// Update modifies an existing hotel document in the Solr collection
func (searchEngine Solr) Update(ctx context.Context, hotel hotels.Hotel) error {
	// Prepare the document for Solr
	doc := map[string]interface{}{
		"id":        hotel.ID,
		"name":      hotel.Name,
		"address":   hotel.Address,
		"city":      hotel.City,
		"state":     hotel.State,
		"rating":    hotel.Rating,
		"amenities": hotel.Amenities,
	}

	// Update the document in Solr
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshaling hotel document: %w", err)
	}

	// Update the document in Solr
	if _, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body)); err != nil {
		return fmt.Errorf("error updating hotel: %w", err)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

func (searchEngine Solr) Delete(ctx context.Context, id string) error {
	// Prepare the delete request
	docToDelete := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": id,
		},
	}

	// Update the document in Solr
	body, err := json.Marshal(docToDelete)
	if err != nil {
		return fmt.Errorf("error marshaling hotel document: %w", err)
	}

	// Execute the delete request using the Update method
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error deleting hotel: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to index hotel: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

func (searchEngine Solr) Search(ctx context.Context, query string) ([]hotels.Hotel, error) {
	// Prepare the Solr query
	solrQuery := fmt.Sprintf("q=%s", query) // Format the query string for Solr

	// Execute the search request
	resp, err := searchEngine.Client.Query(ctx, searchEngine.Collection, solr.NewQuery(solrQuery))
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %w", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("failed to execute search query: %v", resp.Error)
	}

	// Parse the response and extract hotel documents
	var hotelsList []hotels.Hotel
	for _, doc := range resp.Response.Documents {
		// Parse amenities
		amenities := make([]string, 0)
		for _, amenity := range doc["amenities"].([]interface{}) {
			amenities = append(amenities, amenity.(string))
		}

		// Create a hotel from the document fields
		hotel := hotels.Hotel{
			ID:        doc["id"].(string),
			Name:      doc["name"].(string),
			Address:   doc["address"].(string),
			City:      doc["city"].(string),
			State:     doc["state"].(string),
			Rating:    doc["rating"].(float64),
			Amenities: amenities,
		}
		hotelsList = append(hotelsList, hotel)
	}

	return hotelsList, nil
}
