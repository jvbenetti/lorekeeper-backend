package response

import "gorm.io/datatypes"

type EntityResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Lore      string            `json:"lore"`
	ImageURL  string            `json:"image_url"`
	Stats     datatypes.JSON    `json:"stats"`
	Locations []MinimalResource `json:"locations"` // Apenas resumos
	Drops     []MinimalResource `json:"drops"`     // Apenas resumos
}
