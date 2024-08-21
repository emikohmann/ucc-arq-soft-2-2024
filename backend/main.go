// Escribir las líneas del router que definen los endpoints RESTful de un sistema de reserva de habitaciones para hoteles, teniendo en cuenta las siguientes posibles operaciones:
// - Listar todos los hoteles disponibles
// - Buscar hoteles disponibles especificando ciudad de destino y fechas
// - Mostrar el detalle de un hotel seleccionado
// - Listar las reservas de un usuario
// - Hacer una reserva en un hotel para una fecha específica
package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var DoNothing = func(c *gin.Context) {}

type HotelDAO struct {
	ID             int       `bson:"id" json:"id"`
	Name           string    `bson:"name" json:"name"`
	Address        string    `bson:"address" json:"address"`
	City           string    `bson:"city" json:"city"`
	State          string    `bson:"state" json:"state"`
	Country        string    `bson:"country" json:"country"`
	ZipCode        string    `bson:"zip_code" json:"zip_code"`
	PhoneNumber    string    `bson:"phone_number" json:"phone_number"`
	Email          string    `bson:"email" json:"email"`
	Rating         float64   `bson:"rating" json:"rating"`
	AvailableRooms int       `bson:"available_rooms" json:"available_rooms"`
	PricePerNight  float64   `bson:"price_per_night" json:"price_per_night"`
	Amenities      []string  `bson:"amenities" json:"amenities"`
	CheckInTime    time.Time `bson:"check_in_time" json:"check_in_time"`
	CheckOutTime   time.Time `bson:"check_out_time" json:"check_out_time"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}

func main() {
	// Crear un contexto con un tiempo de espera de 10 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configurar las opciones del cliente con URI y credenciales
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		Username: "root",
		Password: "root",
	})

	// Conectarse a MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("Successfully connected to MongoDB")

	// Seleccionar la base de datos y colección
	db := client.Database("testdb")       // Cambia "myDatabase" por tu base de datos
	collection := db.Collection("hotels") // Cambia "users" por tu colección

	// Buscar todos los documentos en la colección
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Guardar los resultados en un slice de bson.M
	// var results []bson.M
	// if err = cursor.All(ctx, &results); err != nil {
	// 	log.Fatal(err)
	// }

	// Guardar los resultados en un slice de Hotel
	var results []HotelDAO
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	// Imprimir los resultados
	for _, hotel := range results {
		fmt.Println(fmt.Sprintf("Hotel '%s' [Address: %s] [Rating: %.2f] [Rooms: %d] [Price per Night: %.2f] [Amenities: %v] [Check-In: %s] [Check-Out: %s] [Created At: %s] [Updated At: %s]",
			hotel.Name,
			hotel.Address,
			hotel.Rating,
			hotel.AvailableRooms,
			hotel.PricePerNight,
			hotel.Amenities,
			hotel.CheckInTime.Format(time.RFC3339),
			hotel.CheckOutTime.Format(time.RFC3339),
			hotel.CreatedAt.Format(time.RFC3339),
			hotel.UpdatedAt.Format(time.RFC3339),
		))
	}

	router := gin.New()
	router.GET("/hotels", DoNothing)             // 1, 2
	router.GET("/hotels/:id", DoNothing)         // 3
	router.GET("/users/:id/bookings", DoNothing) // 4
	router.POST("/bookings", DoNothing)
	router.Run(":8080")

}
