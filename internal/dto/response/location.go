package response

type LocationResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Lore        string            `json:"lore"`
	ImageURL    string            `json:"image_url"`
	SubAreas    []MinimalResource `json:"sub_areas"`
	Entities    []MinimalResource `json:"entities"`
	HiddenItems []MinimalResource `json:"hidden_items"` // Hidden itens (chests)
	EntityDrops []MinimalResource `json:"entity_drops"` // Itens that drop from monster
}
