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

type EntityHandler struct {
	Repo *repository.EntityRepository // Repository Injection
}

// Create handle with create route
func (h *EntityHandler) Create(c echo.Context) error {
	var entity models.Entity

	if err := c.Bind(&entity); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Formato de dados inválido"})
	}

	if !entity.Type.IsValid() {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tipo de local entidade inválido. Tipos permitidos: NPC, GOD, MONSTER, DEMIGOD",
		})
	}

	if entity.Type == models.EntityMonster && len(entity.Stats) > 0 {
		var stats models.MonsterStats
		if err := json.Unmarshal(entity.Stats, &stats); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "O formato dos atributos (stats) do monstro está incorreto",
			})
		}
	}

	// Call repo to save
	if err := h.Repo.Create(&entity); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar entidade"})
	}

	return c.JSON(http.StatusCreated, entity)
}

// AddDrop lida com a rota POST /entities/:id/drops/:item_id
func (h *EntityHandler) AddDrop(c echo.Context) error {
	entID, _ := strconv.Atoi(c.Param("id"))
	itemID, _ := strconv.Atoi(c.Param("item_id"))

	if err := h.Repo.AddDrop(uint(entID), uint(itemID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao vincular drop ao monstro"})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Item adicionado aos drops do monstro com sucesso",
	})
}

// GetByID handle with get route /entities/:id
func (h *EntityHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número"})
	}

	entity, err := h.Repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Entidade não encontrada"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Montando a resposta enxuta
	resp := response.EntityResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		Type:      string(entity.Type),
		Lore:      entity.Lore,
		ImageURL:  entity.ImageURL,
		Stats:     entity.Stats,
		Locations: []response.MinimalResource{},
		Drops:     []response.MinimalResource{},
	}

	for _, loc := range entity.Locations {
		resp.Locations = append(
			resp.Locations,
			response.MinimalResource{
				ID:   loc.ID,
				Name: loc.Name,
				Type: string(loc.Type)})
	}
	for _, drop := range entity.Drops {
		resp.Drops = append(
			resp.Drops,
			response.MinimalResource{
				ID:   drop.ID,
				Name: drop.Name,
				Type: string(drop.Type)})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *EntityHandler) GetAll(c echo.Context) error {
	entities, err := h.Repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao buscar entidades"})
	}

	var resp []response.MinimalResource

	for _, ent := range entities {
		resp = append(resp, response.MinimalResource{
			ID:   ent.ID,
			Name: ent.Name,
			Type: string(ent.Type),
		})
	}

	return c.JSON(http.StatusOK, resp)
}
