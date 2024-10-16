package users

import domain "users-api/domain/users"

type Mock struct{}

func NewMock() Mock {
	return Mock{}
}

func (service Mock) GetAll() ([]domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (service Mock) GetByID(id int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (service Mock) Create(user domain.User) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (service Mock) Login(username string, password string) (domain.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}
