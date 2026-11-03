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

// AddEntity vinc a Monster/NPC that exists and a Local that exists
func (r *LocationRepository) AddEntity(locationID uint, entityID uint) error {
	location := models.Location{ID: locationID}
	entity := models.Entity{ID: entityID}

	// GORM going to in 'location_entities' and insert (location_id, entity_id)
	return r.DB.Model(&location).Association("Entities").Append(&entity)
}

// AddItem vinc a Item that exists and a Local that exists
func (r *LocationRepository) AddItem(locationID uint, itemID uint) error {
	location := models.Location{ID: locationID}
	item := models.Item{ID: itemID}

	// GORM going to in 'location_items' and insert (location_id, item_id)
	return r.DB.Model(&location).Association("Items").Append(&item)
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

func (r *LocationRepository) RemoveItem(locationID uint, itemID uint) error {
	location := models.Location{ID: locationID}
	item := models.Item{ID: itemID}
	return r.DB.Model(&location).Association("Items").Delete(&item)
}
