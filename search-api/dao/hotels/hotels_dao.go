package hotels

type Hotel struct {
	ID        int64    `bson:"id"`
	Name      string   `bson:"name"`
	Address   string   `bson:"address"`
	City      string   `bson:"city"`
	State     string   `bson:"state"`
	Rating    float64  `bson:"rating"`
	Amenities []string `bson:"amenities"`
}
