package request

import (
	"encoding/json"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
)

type CreateLocationRequest struct {
	Name     string              `json:"name"`
	Type     models.LocationType `json:"type"`
	Lore     string              `json:"lore"`
	ImageURL string              `json:"image_url"`
	ParentID *uint               `json:"parent_id"`
}

type UpdateLocationRequest struct {
	Name     *string              `json:"name"`
	Type     *models.LocationType `json:"type"`
	Lore     *string              `json:"lore"`
	ImageURL *string              `json:"image_url"`
	ParentID json.RawMessage      `json:"parent_id"`
}
