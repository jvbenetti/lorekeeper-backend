package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Item represents Weapons, Armor, Artifacts, etc.
type Item struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Type       string         `gorm:"size:50;not null" json:"type"`
	Lore       string         `gorm:"type:text" json:"lore"`
	ImageURL   string         `gorm:"size:500" json:"image_url"`
	Properties datatypes.JSON `json:"properties"` // JSONB for damage, rarity, effects

	// Many-to-Many relationships mapped back from Location and Entity
	Locations []Location `gorm:"many2many:location_items;" json:"locations,omitempty"`
	DroppedBy []Entity   `gorm:"many2many:entity_drops;" json:"dropped_by,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
