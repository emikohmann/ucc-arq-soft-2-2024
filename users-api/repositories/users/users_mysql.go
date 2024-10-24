package users

import (
	"errors"
	"fmt"
	"log"
	"users-api/dao/users"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type MySQL struct {
	db *gorm.DB
}

var (
	migrate = []interface{}{
		users.User{},
	}
)

func NewMySQL(config MySQLConfig) MySQL {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database)

	// Open connection to MySQL using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %s", err.Error())
	}

	// Automigrate structs to Gorm
	for _, target := range migrate {
		if err := db.AutoMigrate(target); err != nil {
			log.Fatalf("error automigrating structs: %s", err.Error())
		}
	}

	return MySQL{
		db: db,
	}
}

func (repository MySQL) GetAll() ([]users.User, error) {
	var usersList []users.User
	if err := repository.db.Find(&usersList).Error; err != nil {
		return nil, fmt.Errorf("error fetching all users: %w", err)
	}
	return usersList, nil
}

func (repository MySQL) GetByID(id int64) (users.User, error) {
	var user users.User
	if err := repository.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("error fetching user by id: %w", err)
	}
	return user, nil
}

func (repository MySQL) GetByUsername(username string) (users.User, error) {
	var user users.User
	if err := repository.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("error fetching user by username: %w", err)
	}
	return user, nil
}

func (repository MySQL) Create(user users.User) (int64, error) {
	if err := repository.db.Create(&user).Error; err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	return user.ID, nil
}

func (repository MySQL) Update(user users.User) error {
	if err := repository.db.Save(&user).Error; err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (repository MySQL) Delete(id int64) error {
	if err := repository.db.Delete(&users.User{}, id).Error; err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
