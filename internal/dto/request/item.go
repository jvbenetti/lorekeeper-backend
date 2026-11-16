package request

import (
	"encoding/json"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
)

type CreateItemRequest struct {
	Name       string          `json:"name"`
	Type       models.ItemType `json:"type"`
	Lore       string          `json:"lore"`
	ImageURL   string          `json:"image_url"`
	Properties json.RawMessage `json:"properties"`
}

type UpdateItemRequest struct {
	Name       *string          `json:"name"`
	Type       *models.ItemType `json:"type"`
	Lore       *string          `json:"lore"`
	ImageURL   *string          `json:"image_url"`
	Properties json.RawMessage  `json:"properties"`
}
