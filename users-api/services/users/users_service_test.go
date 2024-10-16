package users

import (
	"github.com/stretchr/testify/assert"
	"testing"
	domain "users-api/domain/users"
	repositories "users-api/repositories/users"
)

func TestService_Login(t *testing.T) {
	// Set up dependencies
	repository := repositories.NewMock()
	service := NewService(repository)

	// Run test cases
	for _, testCase := range []struct {
		CaseName         string
		Username         string
		Password         string
		ExpectedResponse domain.LoginResponse
		ExpectedError    error
	}{
		{
			CaseName: "Success",
			Username: "test_user",
			Password: "test_password",
			ExpectedResponse: domain.LoginResponse{
				UserID:   0,
				Username: "",
				Token:    "",
			},
			ExpectedError: nil,
		},
	} {
		response, err := service.Login(testCase.Username, testCase.Password)
		assert.Equal(t, testCase.ExpectedResponse, response)
		assert.Equal(t, testCase.ExpectedError, err)
	}
}
