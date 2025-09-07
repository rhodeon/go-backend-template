package domain

import (
	"time"

	dbpets "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/pets"
)

var NewPet = new(Pet)

type Pet struct {
	Id        int64
	Name      string
	Category  PetCategory
	Status    string
	ImageUrls []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Pet) FromDbPet(dbPet dbpets.Pet, category PetCategory) Pet {
	return Pet{
		Id:   dbPet.Id,
		Name: dbPet.Name,
		Category: PetCategory{
			Id:   category.Id,
			Name: category.Name,
		},
		Status:    dbPet.Status.String,
		ImageUrls: dbPet.ImageUrls,
		CreatedAt: dbPet.CreatedAt,
		UpdatedAt: dbPet.UpdatedAt,
	}
}
