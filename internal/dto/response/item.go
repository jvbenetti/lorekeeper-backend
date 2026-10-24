package response

import "gorm.io/datatypes"

type ItemResponse struct {
	ID         uint              `json:"id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Lore       string            `json:"lore"`
	ImageURL   string            `json:"image_url"`
	Properties datatypes.JSON    `json:"properties"`
	Locations  []MinimalResource `json:"locations"`
	DroppedBy  []MinimalResource `json:"dropped_by"`
}
