package repository

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *ItemRepository) GetAll(itemType string, search string, rarity string) ([]models.Item, error) {
	var items []models.Item

	query := r.DB.Model(&models.Item{})

	if itemType != "" {
		query = query.Where("type = ?", itemType)
	}

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if rarity != "" {
		query = query.Where("properties->>'rarity' ILIKE ?", rarity)
	}

	err := query.Find(&items).Error
	return items, err
}

func (r *ItemRepository) Update(item *models.Item) error {
	return r.DB.Omit(clause.Associations).Save(item).Error
}

func (r *ItemRepository) Delete(id uint) error {
	res := r.DB.Delete(&models.Item{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
