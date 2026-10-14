package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EntityType string

const (
	EntityMonster EntityType = "MONSTER"
	EntityNPC     EntityType = "NPC"
	EntityGod     EntityType = "GOD"
	EntityDemigod EntityType = "DEMIGOD"
)

func (e EntityType) IsValid() bool {
	switch e {
	case EntityMonster, EntityNPC, EntityGod, EntityDemigod:
		return true
	}
	return false
}

// Entity represents Monsters, NPCs, Gods, and Demigods
type Entity struct {
	ID       uint           `gorm:"primaryKey" json:"id"`
	Name     string         `gorm:"size:255;not null" json:"name"`
	Type     EntityType     `gorm:"size:50;not null" json:"type"`
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
