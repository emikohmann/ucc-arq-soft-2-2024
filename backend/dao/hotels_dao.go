package dao

type HotelDAO struct {
	ID        int64    `json:"id" bson:"id"`
	Name      string   `json:"name" bson:"name"`
	Address   string   `json:"address" bson:"address"`
	City      string   `json:"city" bson:"city"`
	State     string   `json:"state" bson:"state"`
	Rating    float64  `json:"rating" bson:"rating"`
	Amenities []string `json:"amenities" bson:"amenities"`
}
