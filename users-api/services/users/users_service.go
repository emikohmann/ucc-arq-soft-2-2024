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
		return nil, fmt.Errorf("error getting all users: %s", err.Error())
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
	} else {
		fmt.Println(fmt.Sprintf("warning: error getting user by ID from cache: %s", err.Error()))
	}

	// Check in memcached
	user, err = service.memcachedRepository.GetByID(id)
	if err == nil {
		if _, err := service.cacheRepository.Create(user); err != nil {
			fmt.Println(fmt.Sprintf("warning: error caching user after memcached retrieval: %s", err.Error()))
		}
		return service.convertUser(user), nil
	} else {
		fmt.Println(fmt.Sprintf("warning: error getting user by ID from memcached: %s", err.Error()))
	}

	// Check in main repository
	user, err = service.mainRepository.GetByID(id)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user by ID: %w", err)
	}

	// Save in cache and memcached
	if _, err := service.cacheRepository.Create(user); err != nil {
		fmt.Println(fmt.Sprintf("warning: error caching user after main retrieval: %s", err.Error()))
	}
	if _, err := service.memcachedRepository.Create(user); err != nil {
		fmt.Println(fmt.Sprintf("warning: error saving user in memcached: %s", err.Error()))
	}

	return service.convertUser(user), nil
}

func (service Service) GetByUsername(username string) (domain.User, error) {
	// Check in cache first
	user, err := service.cacheRepository.GetByUsername(username)
	if err == nil {
		return service.convertUser(user), nil
	} else {
		fmt.Println(fmt.Sprintf("warning: error getting user by username from cache: %s", err.Error()))
	}

	// Check memcached
	user, err = service.memcachedRepository.GetByUsername(username)
	if err == nil {
		if _, err := service.cacheRepository.Create(user); err != nil {
			fmt.Println(fmt.Sprintf("warning: error caching user after memcached retrieval: %s", err.Error()))
		}
		return service.convertUser(user), nil
	} else {
		fmt.Println(fmt.Sprintf("warning: error getting user by username from memcached: %s", err.Error()))
	}

	// Check main repository
	user, err = service.mainRepository.GetByUsername(username)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user by username: %w", err)
	}

	// Save in cache and memcached
	if _, err := service.cacheRepository.Create(user); err != nil {
		fmt.Println(fmt.Sprintf("warning: error caching user after main retrieval: %s", err.Error()))
	}
	if _, err := service.memcachedRepository.Create(user); err != nil {
		fmt.Println(fmt.Sprintf("warning: error saving user in memcached: %s", err.Error()))
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
		return 0, fmt.Errorf("error creating user: %s", err.Error())
	}

	// Add to cache and memcached
	newUser.ID = id
	if _, err := service.cacheRepository.Create(newUser); err != nil {
		fmt.Println(fmt.Sprintf("warning: error caching new user: %s", err.Error()))
	}
	if _, err := service.memcachedRepository.Create(newUser); err != nil {
		fmt.Println(fmt.Sprintf("warning: error saving new user in memcached: %s", err.Error()))
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
			return fmt.Errorf("error retrieving existing user: %s", err.Error())
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
		return fmt.Errorf("error updating user: %s", err.Error())
	}

	// Update in cache and memcached
	updatedUser := dao.User{
		ID:       user.ID,
		Username: user.Username,
		Password: passwordHash,
	}
	if err := service.cacheRepository.Update(updatedUser); err != nil {
		fmt.Println(fmt.Sprintf("warning: error updating user in cache: %s", err.Error()))
	}
	if err := service.memcachedRepository.Update(updatedUser); err != nil {
		fmt.Println(fmt.Sprintf("warning: error updating user in memcached: %s", err.Error()))
	}

	return nil
}

func (service Service) Delete(id int64) error {
	// Delete from main repository
	err := service.mainRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting user: %s", err.Error())
	}

	// Delete from cache and memcached
	if err := service.cacheRepository.Delete(id); err != nil {
		fmt.Println(fmt.Sprintf("warning: error deleting user from cache: %s", err.Error()))
	}
	if err := service.memcachedRepository.Delete(id); err != nil {
		fmt.Println(fmt.Sprintf("warning: error deleting user from memcached: %s", err.Error()))
	}

	return nil
}

func (service Service) Login(username string, password string) (domain.LoginResponse, error) {
	// Hash the password
	passwordHash := Hash(password)

	// Try to get user from cache repository first
	user, err := service.cacheRepository.GetByUsername(username)
	if err != nil {
		fmt.Println(fmt.Sprintf("warning: error getting user from cache repository: %s", err.Error()))

		// If not found in cache, try to get user from memcached repository
		user, err = service.memcachedRepository.GetByUsername(username)
		if err != nil {
			fmt.Println(fmt.Sprintf("warning: error getting user from memcached repository: %s", err.Error()))

			// If not found in memcached, try to get user from the main repository (database)
			user, err = service.mainRepository.GetByUsername(username)
			if err != nil {
				// If the user is not found in the main repository, return an error
				return domain.LoginResponse{}, fmt.Errorf("error getting user by username from main repository: %w", err)
			}

			// Save the found user in both cache and memcached repositories
			if _, err := service.cacheRepository.Create(user); err != nil {
				fmt.Println(fmt.Sprintf("warning: error caching user in cache repository: %s", err.Error()))
			}
			if _, err := service.memcachedRepository.Create(user); err != nil {
				fmt.Println(fmt.Sprintf("warning: error caching user in memcached repository: %s", err.Error()))
			}
		} else {
			// Save the found user in the cache repository for future access
			if _, err := service.cacheRepository.Create(user); err != nil {
				fmt.Println(fmt.Sprintf("warning: error caching user in cache repository: %s", err.Error()))
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
