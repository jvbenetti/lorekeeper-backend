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

	// O GORM vai na tabela 'location_entities' e insere (location_id, entity_id)
	return r.DB.Model(&location).Association("Entities").Append(&entity)
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
