package repository

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// RemoveEntity remove the relation between a Location and an Entity
func (r *LocationRepository) RemoveEntity(locationID uint, entityID uint) error {
	location := models.Location{ID: locationID}
	entity := models.Entity{ID: entityID}

	return r.DB.Model(&location).Association("Entities").Delete(&entity)
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

func (r *LocationRepository) GetAll(locationType string, search string) ([]models.Location, error) {
	var locations []models.Location

	query := r.DB.Model(&models.Location{})

	if locationType != "" {
		query = query.Where("type = ?", locationType)
	}

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Find(&locations).Error
	return locations, err
}

func (r *LocationRepository) Update(location *models.Location) error {
	return r.DB.Omit(clause.Associations).Save(location).Error
}

func (r *LocationRepository) Delete(id uint) error {
	res := r.DB.Delete(&models.Location{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *LocationRepository) RemoveItem(locationID uint, itemID uint) error {
	location := models.Location{ID: locationID}
	item := models.Item{ID: itemID}
	return r.DB.Model(&location).Association("Items").Delete(&item)
}
