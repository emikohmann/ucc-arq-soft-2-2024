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
	mainRepository      Repository
	cacheRepository     Repository
	memcachedRepository Repository
	tokenizer           Tokenizer
}

func NewService(mainRepository, cacheRepository, memcachedRepository Repository, tokenizer Tokenizer) Service {
	return Service{
		mainRepository:      mainRepository,
		cacheRepository:     cacheRepository,
		memcachedRepository: memcachedRepository,
		tokenizer:           tokenizer,
	}
}

func (service Service) GetAll() ([]domain.User, error) {
	users, err := service.mainRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	result := make([]domain.User, 0)
	for _, user := range users {
		result = append(result, domain.User{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
		})
	}

	return result, nil
}

func (service Service) GetByID(id int64) (domain.User, error) {
	// Check in cache first
	user, err := service.cacheRepository.GetByID(id)
	if err == nil {
		return service.convertUser(user), nil
	}

	// Check in memcached
	user, err = service.memcachedRepository.GetByID(id)
	if err == nil {
		if _, err := service.cacheRepository.Create(user); err != nil {
			return domain.User{}, fmt.Errorf("error caching user after memcached retrieval: %w", err)
		}
		return service.convertUser(user), nil
	}

	// Check in main repository
	user, err = service.mainRepository.GetByID(id)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user by ID: %w", err)
	}

	// Save in cache and memcached
	if _, err := service.cacheRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error caching user after main retrieval: %w", err)
	}
	if _, err := service.memcachedRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error saving user in memcached: %w", err)
	}

	return service.convertUser(user), nil
}

func (service Service) GetByUsername(username string) (domain.User, error) {
	// Check in cache first
	user, err := service.cacheRepository.GetByUsername(username)
	if err == nil {
		return service.convertUser(user), nil
	}

	// Check memcached
	user, err = service.memcachedRepository.GetByUsername(username)
	if err == nil {
		if _, err := service.cacheRepository.Create(user); err != nil {
			return domain.User{}, fmt.Errorf("error caching user after memcached retrieval: %w", err)
		}
		return service.convertUser(user), nil
	}

	// Check main repository
	user, err = service.mainRepository.GetByUsername(username)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user by username: %w", err)
	}

	// Save in cache and memcached
	if _, err := service.cacheRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error caching user after main retrieval: %w", err)
	}
	if _, err := service.memcachedRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error saving user in memcached: %w", err)
	}

	return service.convertUser(user), nil
}

func (service Service) Create(user domain.User) (int64, error) {
	// Hash the password
	passwordHash := Hash(user.Password)

	newUser := dao.User{
		Username: user.Username,
		Password: passwordHash,
	}

	// Create in main repository
	id, err := service.mainRepository.Create(newUser)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}

	// Add to cache and memcached
	newUser.ID = id
	if _, err := service.cacheRepository.Create(newUser); err != nil {
		return 0, fmt.Errorf("error caching new user: %w", err)
	}
	if _, err := service.memcachedRepository.Create(newUser); err != nil {
		return 0, fmt.Errorf("error saving new user in memcached: %w", err)
	}

	return id, nil
}

func (service Service) Update(user domain.User) error {
	// Hash the password if provided
	var passwordHash string
	if user.Password != "" {
		passwordHash = Hash(user.Password)
	} else {
		existingUser, err := service.mainRepository.GetByID(user.ID)
		if err != nil {
			return fmt.Errorf("error retrieving existing user: %w", err)
		}
		passwordHash = existingUser.Password
	}

	// Update in main repository
	err := service.mainRepository.Update(dao.User{
		ID:       user.ID,
		Username: user.Username,
		Password: passwordHash,
	})
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	// Update in cache and memcached
	updatedUser := dao.User{
		ID:       user.ID,
		Username: user.Username,
		Password: passwordHash,
	}
	if err := service.cacheRepository.Update(updatedUser); err != nil {
		return fmt.Errorf("error updating user in cache: %w", err)
	}
	if err := service.memcachedRepository.Update(updatedUser); err != nil {
		return fmt.Errorf("error updating user in memcached: %w", err)
	}

	return nil
}

func (service Service) Delete(id int64) error {
	// Delete from main repository
	err := service.mainRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	// Delete from cache and memcached
	if err := service.cacheRepository.Delete(id); err != nil {
		return fmt.Errorf("error deleting user from cache: %w", err)
	}
	if err := service.memcachedRepository.Delete(id); err != nil {
		return fmt.Errorf("error deleting user from memcached: %w", err)
	}

	return nil
}

func (service Service) Login(username string, password string) (domain.LoginResponse, error) {
	// Hash the password
	passwordHash := Hash(password)

	// Try to get user from cache repository first
	user, err := service.cacheRepository.GetByUsername(username)
	if err != nil {
		// If not found in cache, try to get user from memcached repository
		user, err = service.memcachedRepository.GetByUsername(username)
		if err != nil {
			// If not found in memcached, try to get user from the main repository (database)
			user, err = service.mainRepository.GetByUsername(username)
			if err != nil {
				return domain.LoginResponse{}, fmt.Errorf("error getting user by username from main repository: %w", err)
			}

			// Save the found user in both cache and memcached repositories
			if _, err := service.cacheRepository.Create(user); err != nil {
				return domain.LoginResponse{}, fmt.Errorf("error caching user in cache repository: %w", err)
			}
			if _, err := service.memcachedRepository.Create(user); err != nil {
				return domain.LoginResponse{}, fmt.Errorf("error caching user in memcached repository: %w", err)
			}
		} else {
			// Save the found user in the cache repository for future access
			if _, err := service.cacheRepository.Create(user); err != nil {
				return domain.LoginResponse{}, fmt.Errorf("error caching user in cache repository: %w", err)
			}
		}
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

	// Send the login response
	return domain.LoginResponse{
		UserID:   user.ID,
		Username: user.Username,
		Token:    token,
	}, nil
}

func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (service Service) convertUser(user dao.User) domain.User {
	return domain.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}
}
