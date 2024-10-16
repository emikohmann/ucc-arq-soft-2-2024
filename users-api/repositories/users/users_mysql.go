package users

import (
	"users-api/dao/users"
)

type MySQLConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type MySQL struct {
}

func NewMySQL(config MySQLConfig) MySQL {
	return MySQL{}
}

func (repository MySQL) GetAll() ([]users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository MySQL) GetByID(id int64) (users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository MySQL) GetByUsername(username string) (users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository MySQL) Create(user users.User) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (repository MySQL) Update(user users.User) error {
	//TODO implement me
	panic("implement me")
}

func (repository MySQL) Delete(id int64) error {
	//TODO implement me
	panic("implement me")
}
