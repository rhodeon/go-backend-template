package responses

import "github.com/rhodeon/go-backend-template/domain"

var NewPet = new(Pet)

type Pet struct {
	Id        int64       `json:"id"`
	Name      string      `json:"name"`
	Category  PetCategory `json:"category"`
	PhotoUrls []string    `json:"photo_urls"`
	Tags      []PetTag    `json:"tags"`
	Status    string      `json:"status"`
}

type PetCategory struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type PetTag struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (u *Pet) FromDomainPet(domainPet domain.Pet) Pet {
	return Pet{
		Id:   domainPet.Id,
		Name: domainPet.Name,
		Category: PetCategory{
			Id:   domainPet.Category.Id,
			Name: domainPet.Category.Name,
		},
		PhotoUrls: domainPet.ImageUrls,
		Tags:      nil,
		Status:    domainPet.Status,
	}
}
