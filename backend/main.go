package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var DoNothing = func(c *gin.Context) {}

type HotelDAO struct {
	ID        int      `json:"id" bson:"id"`
	Name      string   `json:"name" bson:"name"`
	Address   string   `json:"address" bson:"address"`
	City      string   `json:"city" bson:"city"`
	State     string   `json:"state" bson:"state"`
	Rating    float64  `json:"rating" bson:"rating"`
	Amenities []string `json:"amenities" bson:"amenities"`
}

func main() {

	// Creamos un contexto
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second)
	defer cancel()

	// Creamos las opciones
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{
			Username: "root",
			Password: "root",
		})

	// Conectamos el cliente
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Listamos las bases de datos
	result, err := client.ListDatabases(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println(result)

	// router := gin.New()
	// router.GET("/hotels", DoNothing)             // 1, 2
	// router.GET("/hotels/:id", DoNothing)         // 3
	// router.GET("/users/:id/bookings", DoNothing) // 4
	// router.POST("/bookings", DoNothing)
	// router.Run(":8080")
}
