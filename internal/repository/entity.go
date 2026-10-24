package repository

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"gorm.io/gorm"
)

// EntityRepository manages Monsters, NPCs, Gods, etc.
type EntityRepository struct {
	DB *gorm.DB // Receive unique connection pointer of main
}

// Create saves a new entity in the database
func (r *EntityRepository) Create(entity *models.Entity) error {
	return r.DB.Create(entity).Error
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

// AddDrop vinc a Item that exists and a Entity that exists
func (r *EntityRepository) AddDrop(entityID uint, itemID uint) error {
	entity := models.Entity{ID: entityID}
	item := models.Item{ID: itemID}

	// O GORM vai na tabela invisível 'entity_drops' e insere a relação (entity_id, item_id)
	return r.DB.Model(&entity).Association("Drops").Append(&item)
}

// GetByID retrieves an Entity by its ID, preloading many-to-many relationships only if needed.
func (r *EntityRepository) GetByID(id uint) (*models.Entity, error) {
	var entity models.Entity

	// The preload gonna in tables to get local and itens
	err := r.DB.Preload("Locations").Preload("Drops").First(&entity, id).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *EntityRepository) GetAll() ([]models.Entity, error) {
	var entities []models.Entity
	err := r.DB.Find(&entities).Error
	return entities, err
}
