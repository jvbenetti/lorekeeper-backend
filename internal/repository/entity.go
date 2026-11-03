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
