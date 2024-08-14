// Escribir las líneas del router que definen los endpoints RESTful de un sistema de reserva de habitaciones para hoteles, teniendo en cuenta las siguientes posibles operaciones:
// - Listar todos los hoteles disponibles
// - Buscar hoteles disponibles especificando ciudad de destino y fechas
// - Mostrar el detalle de un hotel seleccionado
// - Listar las reservas de un usuario
// - Hacer una reserva en un hotel para una fecha específica
package main

import "github.com/gin-gonic/gin"

var DoNothing = func(c *gin.Context) {}

func main() {
	router := gin.New()
	router.GET("/hotels", DoNothing)             // 1, 2
	router.GET("/hotels/:id", DoNothing)         // 3
	router.GET("/users/:id/bookings", DoNothing) // 4
	router.POST("/bookings", DoNothing)
	router.Run(":8080")
}
