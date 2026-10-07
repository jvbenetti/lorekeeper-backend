package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Entity represents Monsters, NPCs, Gods, and Demigods
type Entity struct {
	ID       uint           `gorm:"primaryKey" json:"id"`
	Name     string         `gorm:"size:255;not null" json:"name"`
	Type     string         `gorm:"size:50;not null" json:"type"`
	Lore     string         `gorm:"type:text" json:"lore"`
	ImageURL string         `gorm:"size:500" json:"image_url"`
	Stats    datatypes.JSON `json:"stats"` // JSONB for dynamic D&D sheets

	// Many-to-Many relationships
	Locations []Location `gorm:"many2many:location_entities;" json:"locations,omitempty"` // Habitats
	Drops     []Item     `gorm:"many2many:entity_drops;" json:"drops,omitempty"`          // Items dropped/owned

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
