package hotels

type Hotel struct {
	ID        string   `bson:"_id,omitempty"`
	Name      string   `bson:"name"`
	Address   string   `bson:"address"`
	City      string   `bson:"city"`
	State     string   `bson:"state"`
	Rating    float64  `bson:"rating"`
	Amenities []string `bson:"amenities"`
}
