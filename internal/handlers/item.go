package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/dto/request"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/dto/response"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ItemHandler struct {
	Repo *repository.ItemRepository
}

func parseAndSanitizeProperties(raw json.RawMessage) (datatypes.JSON, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}

	var properties models.ItemProperties
	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&properties); err != nil {
		return nil, err
	}

	cleanBytes, err := json.Marshal(properties)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(cleanBytes), nil
}

func (h *ItemHandler) Create(c echo.Context) error {
	var req request.CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dados inválidos"})
	}

	name, err := request.ValidateName(req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if !req.Type.IsValid() {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tipo de item inválido. Tipos permitidos: WEAPON, ARMOR, ARTIFACT, POTION",
		})
	}

	if err := request.ValidateLore(req.Lore); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	imageURL, err := request.ValidateImageURL(req.ImageURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	cleanProps, err := parseAndSanitizeProperties(req.Properties)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "O formato das propriedades (properties) do item está incorreto",
		})
	}

	item := models.Item{
		Name:       name,
		Type:       req.Type,
		Lore:       req.Lore,
		ImageURL:   imageURL,
		Properties: cleanProps,
	}

	if err := h.Repo.Create(&item); err != nil {
		c.Logger().Errorf("Erro ao salvar item: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar"})
	}
	return c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número válido"})
	}

	item, err := h.Repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item não encontrado"})
		}
		c.Logger().Errorf("Erro ao buscar item %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro interno ao buscar item"})
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
	itemType := strings.TrimSpace(c.QueryParam("type"))
	search := strings.TrimSpace(c.QueryParam("search"))
	rarity := strings.TrimSpace(c.QueryParam("rarity"))

	if len(search) > request.MaxSearchLength || len(rarity) > request.MaxSearchLength {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Termo de busca muito longo"})
	}

	items, err := h.Repo.GetAll(itemType, search, rarity)
	if err != nil {
		c.Logger().Errorf("Erro ao buscar itens: %v", err)
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
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número válido"})
	}

	item, err := h.Repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item não encontrado"})
		}
		c.Logger().Errorf("Erro ao buscar item %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro interno ao buscar item"})
	}

	var req request.UpdateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dados inválidos"})
	}

	if req.Name != nil {
		name, err := request.ValidateName(*req.Name)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		item.Name = name
	}

	if req.Type != nil {
		if !req.Type.IsValid() {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Tipo de item inválido. Tipos permitidos: WEAPON, ARMOR, ARTIFACT, POTION",
			})
		}
		item.Type = *req.Type
	}

	if req.Lore != nil {
		if err := request.ValidateLore(*req.Lore); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		item.Lore = *req.Lore
	}

	if req.ImageURL != nil {
		imageURL, err := request.ValidateImageURL(*req.ImageURL)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		item.ImageURL = imageURL
	}

	if len(req.Properties) > 0 {
		cleanProps, err := parseAndSanitizeProperties(req.Properties)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "O formato das propriedades (properties) do item está incorreto",
			})
		}
		item.Properties = cleanProps
	}

	if err := h.Repo.Update(item); err != nil {
		c.Logger().Errorf("Erro ao atualizar item %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao atualizar item"})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número válido"})
	}

	if err := h.Repo.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item não encontrado"})
		}
		c.Logger().Errorf("Erro ao deletar item %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao deletar item"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Item removido com sucesso"})
}
