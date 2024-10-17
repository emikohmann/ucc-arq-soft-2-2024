package users_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	dao "users-api/dao/users"
	domain "users-api/domain/users"
	"users-api/internal/tokenizers"
	repositories "users-api/repositories/users"
	service "users-api/services/users"
)

var (
	// Create mocks
	mainRepo      = repositories.NewMock()
	cacheRepo     = repositories.NewMock()
	memcachedRepo = repositories.NewMock()
	tokenizer     = tokenizers.NewMock()
	usersService  = service.NewService(mainRepo, cacheRepo, memcachedRepo, tokenizer)
)

func TestService(t *testing.T) {
	t.Run("GetAll - Success", func(t *testing.T) {
		mockUsers := []dao.User{
			{ID: 1, Username: "user1", Password: "password1"},
			{ID: 2, Username: "user2", Password: "password2"},
		}
		mainRepo.On("GetAll").Return(mockUsers, nil).Once()

		result, err := usersService.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "user1", result[0].Username)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetAll - Error", func(t *testing.T) {
		mainRepo.On("GetAll").Return(nil, errors.New("db error")).Once()

		result, err := usersService.GetAll()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "error getting all users: db error", err.Error()) // Assert expectations    // Assert expectations    // Assert expectations

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Success from Cache", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}
		cacheRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()

		result, err := usersService.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "user1", result.Username)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Not Found in Cache, Found in Memcached", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()
		cacheRepo.On("Create", mockUser).Return(int64(1), nil).Once()

		result, err := usersService.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "user1", result.Username)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Not Found in Cache or Memcached, Found in Main Repo", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		mainRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()
		cacheRepo.On("Create", mockUser).Return(int64(1), nil).Once()
		memcachedRepo.On("Create", mockUser).Return(int64(1), nil).Once()

		result, err := usersService.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "user1", result.Username)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Error in Main Repo", func(t *testing.T) {
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		mainRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("db error")).Once()

		result, err := usersService.GetByID(1)

		assert.Error(t, err)
		assert.Equal(t, "error getting user by ID: db error", err.Error())
		assert.Equal(t, domain.User{}, result)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Create - Success", func(t *testing.T) {
		newUser := dao.User{Username: "newuser", Password: service.Hash("password")}
		mainRepo.On("Create", newUser).Return(int64(1), nil).Once()
		newUser.ID = 1
		cacheRepo.On("Create", newUser).Return(int64(1), nil).Once()
		memcachedRepo.On("Create", newUser).Return(int64(1), nil).Once()

		id, err := usersService.Create(domain.User{Username: "newuser", Password: "password"})

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Create - Error", func(t *testing.T) {
		newUser := dao.User{Username: "newuser", Password: service.Hash("password")}
		mainRepo.On("Create", newUser).Return(int64(0), errors.New("db error")).Once()

		id, err := usersService.Create(domain.User{Username: "newuser", Password: "password"})

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Equal(t, "error creating user: db error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Update - Success", func(t *testing.T) {
		updateUser := dao.User{ID: 1, Username: "updateduser", Password: service.Hash("newpassword")}
		mainRepo.On("Update", updateUser).Return(nil).Once()
		cacheRepo.On("Update", updateUser).Return(nil).Once()
		memcachedRepo.On("Update", updateUser).Return(nil).Once()

		userToUpdate := domain.User{ID: 1, Username: "updateduser", Password: "newpassword"}
		err := usersService.Update(userToUpdate)

		assert.NoError(t, err)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Update - Error", func(t *testing.T) {
		updateUser := dao.User{ID: 1, Username: "updateduser", Password: service.Hash("newpassword")}
		mainRepo.On("Update", updateUser).Return(errors.New("db error")).Once()

		userToUpdate := domain.User{ID: 1, Username: "updateduser", Password: "newpassword"}
		err := usersService.Update(userToUpdate)

		assert.Error(t, err)
		assert.Equal(t, "error updating user: db error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		mainRepo.On("Delete", int64(1)).Return(nil).Once()
		cacheRepo.On("Delete", int64(1)).Return(nil).Once()
		memcachedRepo.On("Delete", int64(1)).Return(nil).Once()

		err := usersService.Delete(1)

		assert.NoError(t, err)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Delete - Error", func(t *testing.T) {
		mainRepo.On("Delete", int64(1)).Return(errors.New("db error")).Once()

		err := usersService.Delete(1)

		assert.Error(t, err)
		assert.Equal(t, "error deleting user: db error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - Success", func(t *testing.T) {
		username := "user1"
		password := "password"
		hashedPassword := service.Hash(password)

		mockUser := dao.User{ID: 1, Username: username, Password: hashedPassword}
		cacheRepo.On("GetByUsername", username).Return(mockUser, nil).Once()
		tokenizer.On("GenerateToken", username, int64(1)).Return("token", nil).Once()

		response, err := usersService.Login(username, password)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), response.UserID)
		assert.Equal(t, "token", response.Token)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - Invalid Credentials", func(t *testing.T) {
		username := "user1"
		password := "wrongpassword"
		hashedPassword := service.Hash("password")

		mockUser := dao.User{ID: 1, Username: username, Password: hashedPassword}
		cacheRepo.On("GetByUsername", username).Return(mockUser, nil).Once()

		response, err := usersService.Login(username, password)

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - User Not Found", func(t *testing.T) {
		username := "user1"
		password := "password"

		cacheRepo.On("GetByUsername", username).Return(dao.User{}, errors.New("not found")).Once()
		memcachedRepo.On("GetByUsername", username).Return(dao.User{}, errors.New("not found")).Once()
		mainRepo.On("GetByUsername", username).Return(dao.User{}, errors.New("not found")).Once()

		response, err := usersService.Login(username, password)

		assert.Error(t, err)
		assert.Equal(t, "error getting user by username from main repository: not found", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - Token Generation Error", func(t *testing.T) {
		username := "user1"
		password := "password"
		hashedPassword := service.Hash(password)

		mockUser := dao.User{ID: 1, Username: username, Password: hashedPassword}
		cacheRepo.On("GetByUsername", username).Return(mockUser, nil).Once()
		tokenizer.On("GenerateToken", username, int64(1)).Return("", errors.New("token error")).Once()

		response, err := usersService.Login(username, password)

		assert.Error(t, err)
		assert.Equal(t, "error generating token: token error", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})
}
