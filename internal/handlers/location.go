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
	"gorm.io/gorm"
)

type LocationHandler struct {
	Repo *repository.LocationRepository
}

func (h *LocationHandler) Create(c echo.Context) error {
	var req request.CreateLocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dados inválidos"})
	}

	name, err := request.ValidateName(req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if !req.Type.IsValid() {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tipo de local inválido. Tipos permitidos: CITY, DUNGEON, HOUSE, CASTLE, HISTORIC",
		})
	}

	if err := request.ValidateLore(req.Lore); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	imageURL, err := request.ValidateImageURL(req.ImageURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.ParentID != nil {
		if *req.ParentID == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do local pai inválido"})
		}
		if _, err := h.Repo.GetByID(*req.ParentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Local pai (parent_id) não encontrado"})
			}
			c.Logger().Errorf("Erro ao validar parent_id %d: %v", *req.ParentID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao validar local pai"})
		}
	}

	location := models.Location{
		Name:     name,
		Type:     req.Type,
		Lore:     req.Lore,
		ImageURL: imageURL,
		ParentID: req.ParentID,
	}

	if err := h.Repo.Create(&location); err != nil {
		c.Logger().Errorf("Erro ao salvar local: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar"})
	}
	return c.JSON(http.StatusCreated, location)
}

// AddEntity handle with POST /locations/:id/entities/:entity_id
func (h *LocationHandler) AddEntity(c echo.Context) error {
	locID, err := strconv.Atoi(c.Param("id"))
	if err != nil || locID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do local deve ser um número válido"})
	}

	entID, err := strconv.Atoi(c.Param("entity_id"))
	if err != nil || entID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID da entidade deve ser um número válido"})
	}

	if err := h.Repo.AddEntity(uint(locID), uint(entID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Local ou entidade não encontrado"})
		}
		c.Logger().Errorf("Erro ao vincular entidade ao local: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao vincular entidade ao local"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Entidade vinculada ao local com sucesso",
	})
}

// AddItem handle with POST /locations/:id/items/:item_id
func (h *LocationHandler) AddItem(c echo.Context) error {
	locID, err := strconv.Atoi(c.Param("id"))
	if err != nil || locID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do local deve ser um número válido"})
	}

	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil || itemID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do item deve ser um número válido"})
	}

	if err := h.Repo.AddItem(uint(locID), uint(itemID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Local ou item não encontrado"})
		}
		c.Logger().Errorf("Erro ao vincular item ao local: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao vincular item ao local"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Item escondido no local com sucesso",
	})
}

// RemoveEntity handle with DELETE /locations/:id/entities/:entity_id
func (h *LocationHandler) RemoveEntity(c echo.Context) error {
	locID, err := strconv.Atoi(c.Param("id"))
	if err != nil || locID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do local deve ser um número válido"})
	}

	entID, err := strconv.Atoi(c.Param("entity_id"))
	if err != nil || entID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID da entidade deve ser um número válido"})
	}

	if err := h.Repo.RemoveEntity(uint(locID), uint(entID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Local ou entidade não encontrado"})
		}
		c.Logger().Errorf("Erro ao remover entidade do local: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao remover entidade do local"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Entidade removida do local com sucesso",
	})
}

// RemoveItem handle with DELETE /locations/:id/items/:item_id
func (h *LocationHandler) RemoveItem(c echo.Context) error {
	locID, err := strconv.Atoi(c.Param("id"))
	if err != nil || locID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do local deve ser um número válido"})
	}

	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil || itemID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do item deve ser um número válido"})
	}

	if err := h.Repo.RemoveItem(uint(locID), uint(itemID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Local ou item não encontrado"})
		}
		c.Logger().Errorf("Erro ao remover item do local: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao remover item do local"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Item removido do local com sucesso",
	})
}

func (h *LocationHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número válido"})
	}

	location, err := h.Repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Local não encontrado"})
		}
		c.Logger().Errorf("Erro ao buscar local %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro interno ao buscar local"})
	}

	resp := response.LocationResponse{
		ID:          location.ID,
		Name:        location.Name,
		Type:        string(location.Type),
		Lore:        location.Lore,
		ImageURL:    location.ImageURL,
		SubAreas:    []response.MinimalResource{},
		Entities:    []response.MinimalResource{},
		HiddenItems: []response.MinimalResource{},
		EntityDrops: []response.MinimalResource{},
	}

	for _, sub := range location.SubAreas {
		resp.SubAreas = append(
			resp.SubAreas,
			response.MinimalResource{
				ID:   sub.ID,
				Name: sub.Name,
				Type: string(sub.Type)})
	}
	for _, item := range location.Items {
		resp.HiddenItems = append(
			resp.HiddenItems,
			response.MinimalResource{
				ID:   item.ID,
				Name: item.Name,
				Type: string(item.Type)})
	}

	// Get monsters drops
	seenDrops := make(map[uint]bool)
	for _, ent := range location.Entities {
		resp.Entities = append(
			resp.Entities, response.MinimalResource{ID: ent.ID, Name: ent.Name, Type: string(ent.Type)})

		for _, drop := range ent.Drops {
			if !seenDrops[drop.ID] {
				seenDrops[drop.ID] = true
				resp.EntityDrops = append(
					resp.EntityDrops,
					response.MinimalResource{
						ID:   drop.ID,
						Name: drop.Name,
						Type: string(drop.Type)})
			}
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *LocationHandler) GetAll(c echo.Context) error {
	locationType := strings.TrimSpace(c.QueryParam("type"))
	search := strings.TrimSpace(c.QueryParam("search"))

	if len(search) > request.MaxSearchLength {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Termo de busca muito longo"})
	}

	locations, err := h.Repo.GetAll(locationType, search)
	if err != nil {
		c.Logger().Errorf("Erro ao buscar locais: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao buscar locais"})
	}

	resp := []response.MinimalResource{}

	for _, loc := range locations {
		resp = append(resp, response.MinimalResource{
			ID:   loc.ID,
			Name: loc.Name,
			Type: string(loc.Type),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *LocationHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número válido"})
	}

	location, err := h.Repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Local não encontrado"})
		}
		c.Logger().Errorf("Erro ao buscar local %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro interno ao buscar local"})
	}

	var req request.UpdateLocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dados inválidos"})
	}

	if req.Name != nil {
		name, err := request.ValidateName(*req.Name)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		location.Name = name
	}

	if req.Type != nil {
		if !req.Type.IsValid() {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Tipo de local inválido. Tipos permitidos: CITY, DUNGEON, HOUSE, CASTLE, HISTORIC",
			})
		}
		location.Type = *req.Type
	}

	if req.Lore != nil {
		if err := request.ValidateLore(*req.Lore); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		location.Lore = *req.Lore
	}

	if req.ImageURL != nil {
		imageURL, err := request.ValidateImageURL(*req.ImageURL)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		location.ImageURL = imageURL
	}

	if len(req.ParentID) > 0 {
		trimmed := bytes.TrimSpace(req.ParentID)
		if bytes.Equal(trimmed, []byte("null")) {
			location.ParentID = nil
		} else {
			var parentID uint
			if err := json.Unmarshal(trimmed, &parentID); err != nil || parentID == 0 {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do local pai inválido"})
			}
			if parentID == uint(id) {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Um local não pode ser sub-área de si mesmo"})
			}
			if _, err := h.Repo.GetByID(parentID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return c.JSON(http.StatusBadRequest, map[string]string{"error": "Local pai (parent_id) não encontrado"})
				}
				c.Logger().Errorf("Erro ao validar parent_id %d: %v", parentID, err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao validar local pai"})
			}
			location.ParentID = &parentID
		}
	}

	if err := h.Repo.Update(location); err != nil {
		c.Logger().Errorf("Erro ao atualizar local %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao atualizar local"})
	}

	return c.JSON(http.StatusOK, location)
}

func (h *LocationHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número válido"})
	}

	if err := h.Repo.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Local não encontrado"})
		}
		c.Logger().Errorf("Erro ao deletar local %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao deletar local"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Local removido com sucesso"})
}
