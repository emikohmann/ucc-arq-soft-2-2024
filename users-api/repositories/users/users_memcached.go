package users

import "users-api/dao/users"

type MemcachedConfig struct {
}

type Memcached struct{}

func NewMemcached(config MemcachedConfig) Memcached {
	return Memcached{}
}

func (repository Memcached) GetAll() ([]users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository Memcached) GetByID(id int64) (users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository Memcached) GetByUsername(username string) (users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository Memcached) Create(user users.User) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (repository Memcached) Update(user users.User) error {
	//TODO implement me
	panic("implement me")
}

func (repository Memcached) Delete(id int64) error {
	//TODO implement me
	panic("implement me")
}
