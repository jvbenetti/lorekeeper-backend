package request

import (
	"encoding/json"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
)

type CreateEntityRequest struct {
	Name     string          `json:"name"`
	Type     models.EntityType `json:"type"`
	Lore     string          `json:"lore"`
	ImageURL string          `json:"image_url"`
	Stats    json.RawMessage `json:"stats"`
}

type UpdateEntityRequest struct {
	Name     *string            `json:"name"`
	Type     *models.EntityType `json:"type"`
	Lore     *string            `json:"lore"`
	ImageURL *string            `json:"image_url"`
	Stats    json.RawMessage    `json:"stats"`
}
