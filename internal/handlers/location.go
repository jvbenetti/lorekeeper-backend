package handlers

import (
	"net/http"
	"strconv"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/dto/response"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/repository"
	"github.com/labstack/echo/v4"
)

type LocationHandler struct {
	Repo *repository.LocationRepository
}

func (h *LocationHandler) Create(c echo.Context) error {
	var location models.Location
	if err := c.Bind(&location); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dados inválidos"})
	}

	if !location.Type.IsValid() {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tipo de local inválido. Tipos permitidos: CITY, DUNGEON, HOUSE, CASTLE, HISTORIC",
		})
	}

	if err := h.Repo.Create(&location); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar"})
	}
	return c.JSON(http.StatusCreated, location)
}

// AddEntity handle with POST /locations/:id/entities/:entity_id
func (h *LocationHandler) AddEntity(c echo.Context) error {
	locID, _ := strconv.Atoi(c.Param("id"))
	entID, _ := strconv.Atoi(c.Param("entity_id"))

	if err := h.Repo.AddEntity(uint(locID), uint(entID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao vincular entidade ao local"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Entidade vinculada ao local com sucesso",
	})
}

// AddItem handle with POST /locations/:id/items/:item_id
func (h *LocationHandler) AddItem(c echo.Context) error {
	locID, _ := strconv.Atoi(c.Param("id"))
	itemID, _ := strconv.Atoi(c.Param("item_id"))

	if err := h.Repo.AddItem(uint(locID), uint(itemID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao vincular item ao local"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Item escondido no local com sucesso",
	})
}

func (h *LocationHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número"})
	}

	location, err := h.Repo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Local não encontrado"})
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
	for _, ent := range location.Entities {
		resp.Entities = append(
			resp.Entities, response.MinimalResource{ID: ent.ID, Name: ent.Name, Type: string(ent.Type)})

		for _, drop := range ent.Drops {
			resp.EntityDrops = append(
				resp.EntityDrops,
				response.MinimalResource{
					ID:   drop.ID,
					Name: drop.Name,
					Type: string(drop.Type)})
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *LocationHandler) GetAll(c echo.Context) error {
	locations, err := h.Repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao buscar locais"})
	}

	var resp []response.MinimalResource

	for _, loc := range locations {
		resp = append(resp, response.MinimalResource{
			ID:   loc.ID,
			Name: loc.Name,
			Type: string(loc.Type),
		})
	}

	return c.JSON(http.StatusOK, resp)
}
