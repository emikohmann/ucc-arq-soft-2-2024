package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	controllers "users-api/controllers/users"
	"users-api/internal/tokenizers"
	repositories "users-api/repositories/users"
	services "users-api/services/users"
)

func main() {
	// MySQL
	mySQLConfig := repositories.MySQLConfig{
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "users-api",
		Username: "root",
		Password: "root",
	}
	mySQLRepo := repositories.NewMySQL(mySQLConfig)

	// JWT
	jwtConfig := tokenizers.JWTConfig{
		Key:      "ThisIsAnExampleJWTKey!",
		Duration: 24 * time.Hour,
	}
	jwtTokenizer := tokenizers.NewTokenizer(jwtConfig)

	// Services
	service := services.NewService(mySQLRepo, jwtTokenizer)

	// Handlers
	controller := controllers.NewController(service)

	// Create router
	router := gin.Default()

	// URL mappings
	router.GET("/users", controller.GetAll)
	router.GET("/users/:id", controller.GetByID)
	router.POST("/users", controller.Create)
	router.POST("/login", controller.Login)

	// Run application
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running application: %v", err)
	}
}
