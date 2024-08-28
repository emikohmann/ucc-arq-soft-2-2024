package main

import (
	"backend/repositories"
	"fmt"
)

func main() {

	var repo repositories.HotelsRepo
	repo = repositories.NewHotelsMongo()
	fmt.Println(repo.GetHotelByID(3))
}
