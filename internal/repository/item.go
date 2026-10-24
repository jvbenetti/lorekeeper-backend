package repository

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"gorm.io/gorm"
)

type ItemRepository struct {
	DB *gorm.DB
}

func (r *ItemRepository) Create(item *models.Item) error {
	return r.DB.Create(item).Error
}

func (r *ItemRepository) GetByID(id uint) (*models.Item, error) {
	var item models.Item
	// Preload get locations where this item can be found and the entities that drop it
	err := r.DB.Preload("Locations").Preload("DroppedBy").First(&item, id).Error
	return &item, err
}

func (r *ItemRepository) GetAll() ([]models.Item, error) {
	var items []models.Item
	err := r.DB.Find(&items).Error
	return items, err
}
