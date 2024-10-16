package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	dao "users-api/dao/users"
	domain "users-api/domain/users"
)

type Repository interface {
	GetAll() ([]dao.User, error)
	GetByID(id int64) (dao.User, error)
	GetByUsername(username string) (dao.User, error)
	Create(user dao.User) (int64, error)
	Update(user dao.User) error
	Delete(id int64) error
}

type Tokenizer interface {
	GenerateToken(username string, userID int64) (string, error)
}

type Service struct {
	repository Repository
	tokenizer  Tokenizer
}

func NewService(repository Repository, tokenizer Tokenizer) Service {
	return Service{
		repository: repository,
		tokenizer:  tokenizer,
	}
}

func (service Service) GetAll() ([]domain.User, error) {
	// Try to get all users
	users, err := service.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %s", err.Error())
	}

	// Convert
	result := make([]domain.User, 0)
	for _, user := range users {
		result = append(result, domain.User{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
		})
	}

	// Send the result
	return result, nil
}

func (service Service) GetByID(id int64) (domain.User, error) {
	// Try to get the user by ID
	user, err := service.repository.GetByID(id)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user by ID: %w", err)
	}

	// Send the user
	return domain.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

func (service Service) Login(username string, password string) (domain.LoginResponse, error) {
	// Hash the password
	passwordHash := Hash(password)

	// Try to get user by username
	user, err := service.repository.GetByUsername(username)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("error getting user by username: %w", err)
	}

	// Compare passwords
	if user.Password != passwordHash {
		return domain.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	// Generate token
	token, err := service.tokenizer.GenerateToken(user.Username, user.ID)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("error generating token: %w", err)
	}

	// Send the login
	return domain.LoginResponse{
		UserID:   user.ID,
		Username: user.Username,
		Token:    token,
	}, nil
}

func (service Service) Create(user domain.User) (int64, error) {
	// Hash the password
	passwordHash := Hash(user.Password)

	// Try to create the user
	id, err := service.repository.Create(dao.User{
		Username: user.Username,
		Password: passwordHash,
	})
	if err != nil {
		return 0, fmt.Errorf("error creating user: %s", err.Error())
	}

	// Send ID
	return id, nil
}

func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
