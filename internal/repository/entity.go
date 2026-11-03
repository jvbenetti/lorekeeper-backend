package repository

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EntityRepository manages Monsters, NPCs, Gods, etc.
type EntityRepository struct {
	DB *gorm.DB // Receive unique connection pointer of main
}

// Create saves a new entity in the database
func (r *EntityRepository) Create(entity *models.Entity) error {
	return r.DB.Create(entity).Error
}

// AddDrop vinc a Item that exists and a Entity that exists
func (r *EntityRepository) AddDrop(entityID uint, itemID uint) error {
	entity := models.Entity{ID: entityID}
	item := models.Item{ID: itemID}

	// O GORM vai na tabela invisível 'entity_drops' e insere a relação (entity_id, item_id)
	return r.DB.Model(&entity).Association("Drops").Append(&item)
}

// RemoveDrop remove the relation between an Entity and an Item drop
func (r *EntityRepository) RemoveDrop(entityID uint, itemID uint) error {
	entity := models.Entity{ID: entityID}
	item := models.Item{ID: itemID}

	return r.DB.Model(&entity).Association("Drops").Delete(&item)
}

// Update updates an existing entity in the database
func (r *EntityRepository) Update(entity *models.Entity) error {
	return r.DB.Omit(clause.Associations).Save(entity).Error
}

// Delete soft-deletes an entity by ID
func (r *EntityRepository) Delete(id uint) error {
	res := r.DB.Delete(&models.Entity{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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

// GetAll get all Entity's, with filters
func (r *EntityRepository) GetAll(entityType string, search string) ([]models.Entity, error) {
	var entities []models.Entity

	query := r.DB.Model(&models.Entity{})

	// If user sent a type
	if entityType != "" {
		query = query.Where("type = ?", entityType)
	}

	// If user sent a search term
	if search != "" {
		// ILIKE in PostgreSQL ignore uppercases and downcases
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// Executa a busca no banco
	err := query.Find(&entities).Error
	return entities, err
}
