package responses

import "time"

type Order struct {
	Id       int       `json:"id" required:"true"`
	PetId    int       `json:"pet_id" required:"true"`
	Quantity int       `json:"quantity" required:"true"`
	ShipDate time.Time `json:"ship_date" required:"true"`
	Status   string    `json:"status" enum:"placed,approved,delivered" required:"true"`
	Complete bool      `json:"complete" required:"true"`
}
