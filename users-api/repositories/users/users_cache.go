package users

import "users-api/dao/users"

type CacheConfig struct {
}

type Cache struct{}

func NewCache(config CacheConfig) Cache {
	return Cache{}
}

func (repository Cache) GetAll() ([]users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository Cache) GetByID(id int64) (users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository Cache) GetByUsername(username string) (users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository Cache) Create(user users.User) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (repository Cache) Update(user users.User) error {
	//TODO implement me
	panic("implement me")
}

func (repository Cache) Delete(id int64) error {
	//TODO implement me
	panic("implement me")
}
