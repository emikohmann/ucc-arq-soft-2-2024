package users

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"users-api/dao/users"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

type MySQLConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type MySQL struct {
	db *sql.DB
}

func NewMySQL(config MySQLConfig) MySQL {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Username, config.Password, config.Host, config.Port, config.Database)

	// Open connection to MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %s", err.Error())
	}

	// Ping the database to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping MySQL: %s", err.Error())
	}

	return MySQL{
		db: db,
	}
}

func (repository MySQL) GetAll() ([]users.User, error) {
	rows, err := repository.db.Query("SELECT id, username, password FROM users")
	if err != nil {
		return nil, fmt.Errorf("error fetching all users: %w", err)
	}
	defer rows.Close()

	var usersList []users.User
	for rows.Next() {
		var user users.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		usersList = append(usersList, user)
	}
	return usersList, nil
}

func (repository MySQL) GetByID(id int64) (users.User, error) {
	var user users.User
	if err := repository.db.
		QueryRow("SELECT id, username, password FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("error fetching user by id: %w", err)
	}
	return user, nil
}

func (repository MySQL) GetByUsername(username string) (users.User, error) {
	var user users.User
	if err := repository.db.
		QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("error fetching user by username: %w", err)
	}
	return user, nil
}

func (repository MySQL) Create(user users.User) (int64, error) {
	result, err := repository.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %w", err)
	}
	return id, nil
}

func (repository MySQL) Update(user users.User) error {
	if _, err := repository.db.Exec("UPDATE users SET username = ?, password = ? WHERE id = ?", user.Username, user.Password, user.ID); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (repository MySQL) Delete(id int64) error {
	if _, err := repository.db.Exec("DELETE FROM users WHERE id = ?", id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
