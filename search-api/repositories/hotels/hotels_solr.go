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
	Host       string // Solr host
	Port       string // Solr port
	Collection string // Solr collection name
}

type Solr struct {
	Client     *solr.JSONClient
	Collection string
}

// NewSolr initializes a new Solr client
func NewSolr(config SolrConfig) Solr {
	// Construct the BaseURL using the provided host and port
	baseURL := fmt.Sprintf("http://%s:%s", config.Host, config.Port)
	client := solr.NewJSONClient(baseURL)

	return Solr{
		Client:     client,
		Collection: config.Collection,
	}
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

	// Prepare the index request
	indexRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Index the document in Solr
	body, err := json.Marshal(indexRequest)
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

	// Prepare the update request
	updateRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Update the document in Solr
	body, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("error marshaling hotel document: %w", err)
	}

	// Execute the update request using the Update method
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error updating hotel: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to update hotel: %v", resp.Error)
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

func (searchEngine Solr) Search(ctx context.Context, query string, limit int, offset int) ([]hotels.Hotel, error) {
	// Prepare the Solr query with limit and offset
	solrQuery := fmt.Sprintf("q=(name:%s)&rows=%d&start=%d", query, limit, offset)

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
		// Initialize amenities slice
		var amenities []string

		// Check if amenities exist and handle different types
		if amenitiesData, ok := doc["amenities"].([]interface{}); ok {
			for _, amenity := range amenitiesData {
				if amenityStr, ok := amenity.(string); ok {
					amenities = append(amenities, amenityStr)
				}
			}
		}

		// Safely extract hotel fields with type assertions
		hotel := hotels.Hotel{
			ID:        getStringField(doc, "id"),
			Name:      getStringField(doc, "name"),
			Address:   getStringField(doc, "address"),
			City:      getStringField(doc, "city"),
			State:     getStringField(doc, "state"),
			Rating:    getFloatField(doc, "rating"),
			Amenities: amenities,
		}
		hotelsList = append(hotelsList, hotel)
	}

	return hotelsList, nil
}

// Helper function to safely get string fields from the document
func getStringField(doc map[string]interface{}, field string) string {
	if val, ok := doc[field].(string); ok {
		return val
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if strVal, ok := val[0].(string); ok {
			return strVal
		}
	}
	return ""
}

// Helper function to safely get float64 fields from the document
func getFloatField(doc map[string]interface{}, field string) float64 {
	if val, ok := doc[field].(float64); ok {
		return val
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if floatVal, ok := val[0].(float64); ok {
			return floatVal
		}
	}
	return 0.0
}
