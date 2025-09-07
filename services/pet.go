package services

import (
	"context"

	"github.com/rhodeon/go-backend-template/domain"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/repositories"
	dbpets "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/pets"

	"github.com/go-errors/errors"
)

type Pet struct {
	*service
}

var petService *Pet

func newPet(repos *repositories.Repositories, cfg *Config) *Pet {
	petService = &Pet{newService(repos, cfg)}
	return petService
}

func (p *Pet) Create(ctx context.Context, dbTx *database.Tx, pet domain.Pet) (domain.Pet, error) {
	createdPet, err := p.repos.Database.Pets.Create(ctx, dbTx, dbpets.CreateParams{
		Name:       pet.Name,
		CategoryId: pet.Category.Id,
		ImageUrls:  pet.ImageUrls,
	})
	if err != nil {
		return domain.Pet{}, errors.Errorf("creating pet in database: %w", err)
	}

	petCategory, err := p.GetCategoryById(ctx, dbTx, pet.Category.Id)
	if err != nil {
		return domain.Pet{}, err
	}

	return domain.NewPet.FromDbPet(createdPet, petCategory), nil
}

func (p *Pet) GetCategoryById(ctx context.Context, dbTx *database.Tx, petCategoryId int64) (domain.PetCategory, error) {
	dbPetCategory, err := p.repos.Database.PetCategories.GetById(ctx, dbTx, petCategoryId)
	if err != nil {
		return domain.PetCategory{}, errors.Errorf("getting pet category by id from database: %w", err)
	}

	return domain.NewPetCategory.FromDbPetCategory(dbPetCategory), nil
}
