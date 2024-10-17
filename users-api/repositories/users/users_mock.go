package users

import (
	"github.com/stretchr/testify/mock"
	"users-api/dao/users"
)

// Mock the Repository and Tokenizer interfaces
type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}
func (m *Mock) GetAll() ([]users.User, error) {
	args := m.Called()
	if err := args.Error(1); err != nil {
		return nil, err // Return early if there's an error
	}
	return args.Get(0).([]users.User), nil
}

func (m *Mock) GetByID(id int64) (users.User, error) {
	args := m.Called(id)
	if err := args.Error(1); err != nil {
		return users.User{}, err // Return zero User if there's an error
	}
	return args.Get(0).(users.User), nil
}

func (m *Mock) GetByUsername(username string) (users.User, error) {
	args := m.Called(username)
	if err := args.Error(1); err != nil {
		return users.User{}, err // Return zero User if there's an error
	}
	return args.Get(0).(users.User), nil
}

func (m *Mock) Create(user users.User) (int64, error) {
	args := m.Called(user)
	if err := args.Error(1); err != nil {
		return 0, err // Return 0 if there's an error
	}
	return args.Get(0).(int64), nil
}

func (m *Mock) Update(user users.User) error {
	args := m.Called(user)
	return args.Error(0) // No change needed here as it returns an error directly
}

func (m *Mock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0) // No change needed here as it returns an error directly
}
