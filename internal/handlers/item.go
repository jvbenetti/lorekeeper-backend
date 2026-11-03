package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/dto/response"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ItemHandler struct {
	Repo *repository.ItemRepository
}

func (h *ItemHandler) Create(c echo.Context) error {
	var item models.Item
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dados inválidos"})
	}

	if !item.Type.IsValid() {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tipo de item inválido. Tipos permitidos: WEAPON, ARMOR, ARTIFACT, POTION",
		})
	}

	if len(item.Properties) > 0 {
		var properties models.ItemProperties
		if err := json.Unmarshal(item.Properties, &properties); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "O formato das propriedades (properties) do item está incorreto",
			})
		}
	}

	if err := h.Repo.Create(&item); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar"})
	}
	return c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número"})
	}

	item, err := h.Repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item não encontrado"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resp := response.ItemResponse{
		ID:         item.ID,
		Name:       item.Name,
		Type:       string(item.Type),
		Lore:       item.Lore,
		ImageURL:   item.ImageURL,
		Properties: item.Properties,
		Locations:  []response.MinimalResource{},
		DroppedBy:  []response.MinimalResource{},
	}

	for _, loc := range item.Locations {
		resp.Locations = append(resp.Locations, response.MinimalResource{
			ID:   loc.ID,
			Name: loc.Name,
			Type: string(loc.Type),
		})
	}

	for _, ent := range item.DroppedBy {
		resp.DroppedBy = append(resp.DroppedBy, response.MinimalResource{
			ID:   ent.ID,
			Name: ent.Name,
			Type: string(ent.Type),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ItemHandler) GetAll(c echo.Context) error {
	itemType := c.QueryParam("type")
	search := c.QueryParam("search")
	rarity := c.QueryParam("rarity")

	items, err := h.Repo.GetAll(itemType, search, rarity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao buscar itens"})
	}

	resp := []response.MinimalResource{}

	for _, ent := range items {
		resp = append(resp, response.MinimalResource{
			ID:   ent.ID,
			Name: ent.Name,
			Type: string(ent.Type),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ItemHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número"})
	}

	item, err := h.Repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item não encontrado"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dados inválidos"})
	}
	item.ID = uint(id)

	if !item.Type.IsValid() {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tipo de item inválido. Tipos permitidos: WEAPON, ARMOR, ARTIFACT, POTION",
		})
	}

	if len(item.Properties) > 0 {
		var properties models.ItemProperties
		if err := json.Unmarshal(item.Properties, &properties); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "O formato das propriedades (properties) do item está incorreto",
			})
		}
	}

	if err := h.Repo.Update(item); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao atualizar item"})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número"})
	}

	if err := h.Repo.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item não encontrado"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao deletar item"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Item removido com sucesso"})
}
