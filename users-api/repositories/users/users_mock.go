package users

import (
	"fmt"
	"users-api/dao/users"
)

type Mock struct {
	data map[int64]users.User
}

func NewMock() Mock {
	return Mock{
		data: make(map[int64]users.User),
	}
}

func (repository Mock) GetAll() ([]users.User, error) {
	result := make([]users.User, 0)
	for _, user := range repository.data {
		result = append(result, user)
	}
	return result, nil
}

func (repository Mock) GetByID(id int64) (users.User, error) {
	user, exists := repository.data[id]
	if !exists {
		return users.User{}, fmt.Errorf("not found")
	}
	return user, nil
}

func (repository Mock) GetByUsername(username string) (users.User, error) {
	for _, user := range repository.data {
		if user.Username == username {
			return user, nil
		}
	}
	return users.User{}, fmt.Errorf("not found")
}

func (repository Mock) Create(user users.User) (int64, error) {
	for _, existingUser := range repository.data {
		if existingUser.Username == user.Username {
			return 0, fmt.Errorf("already exists")
		}
	}

	id := int64(len(repository.data) + 1)
	repository.data[id] = users.User{
		ID:       id,
		Username: user.Username,
		Password: user.Password,
	}

	return id, nil
}

func (repository Mock) Update(user users.User) error {
	repository.data[user.ID] = user
	return nil
}

func (repository Mock) Delete(id int64) error {
	delete(repository.data, id)
	return nil
}
