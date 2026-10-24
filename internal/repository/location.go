package repository

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"gorm.io/gorm"
)

type LocationRepository struct {
	DB *gorm.DB
}

func (r *LocationRepository) Create(location *models.Location) error {
	return r.DB.Create(location).Error
}

func (r *LocationRepository) GetByID(id uint) (*models.Location, error) {
	var location models.Location
	// Preload get itens and monsters of location
	err := r.DB.
		Preload("Entities").
		Preload("Items").
		Preload("SubAreas").
		Preload("Entities.Drops").
		First(&location, id).Error

	return &location, err
}

func (r *LocationRepository) GetAll() ([]models.Location, error) {
	var locations []models.Location
	err := r.DB.Find(&locations).Error
	return locations, err
}
