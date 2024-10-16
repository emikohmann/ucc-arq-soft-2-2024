package search

import (
	"context"
	"github.com/stretchr/testify/assert"
	"search-api/dao/hotels"
	repositories "search-api/repositories/hotels"
	"testing"
)

func TestService_Search(t *testing.T) {
	// Set up dependencies
	ctx := context.Background()
	repository := repositories.NewMock()
	service := NewService(repository)

	// Run test cases
	for _, testCase := range []struct {
		CaseName         string
		Query            string
		Offset           int
		Limit            int
		ExpectedResponse []hotels.Hotel
		ExpectedError    error
	}{
		{
			CaseName:         "Success",
			Query:            "test query",
			Offset:           0,
			Limit:            10,
			ExpectedResponse: []hotels.Hotel{},
			ExpectedError:    nil,
		},
	} {
		response, err := service.Search(ctx, testCase.Query, testCase.Offset, testCase.Limit)
		assert.Equal(t, testCase.ExpectedResponse, response)
		assert.Equal(t, testCase.ExpectedError, err)
	}
}
