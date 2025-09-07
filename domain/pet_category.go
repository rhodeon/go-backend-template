package domain

import (
	dbpetcategories "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/petcategories"
)

var NewPetCategory = new(PetCategory)

type PetCategory struct {
	Id   int64
	Name string
}

func (pc PetCategory) FromDbPetCategory(dbPetCategory dbpetcategories.PetCategory) PetCategory {
	return PetCategory{
		Id:   dbPetCategory.Id,
		Name: dbPetCategory.Name,
	}
}
