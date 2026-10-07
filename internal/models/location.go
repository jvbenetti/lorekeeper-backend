package models

import (
	"time"

	"gorm.io/gorm"
)

// Location represents Cities, Dungeons, Cemeteries, etc.
type Location struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"size:255;not null" json:"name"`
	Type     string `gorm:"size:50;not null" json:"type"`
	Lore     string `gorm:"type:text" json:"lore"`
	ImageURL string `gorm:"size:500" json:"image_url"`

	// Hierarchy (A dungeon inside a city, etc.)
	ParentID *uint      `json:"parent_id"`
	Parent   *Location  `gorm:"foreignKey:ParentID" json:"-"`
	SubAreas []Location `gorm:"foreignKey:ParentID" json:"sub_areas,omitempty"`

	// Many-to-Many relationships
	Entities []Entity `gorm:"many2many:location_entities;" json:"entities,omitempty"` // Monsters/NPCs found here
	Items    []Item   `gorm:"many2many:location_items;" json:"items,omitempty"`       // Loot/Items hidden here

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
