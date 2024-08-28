package main

import (
	"fmt"
	"hotels-api/repositories/hotels"
)

func main() {
	var repo hotels.HotelsRepo
	repo = hotels.NewHotelsMongo()
	fmt.Println(repo.GetHotelByID(3))
}
