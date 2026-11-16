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

type EntityHandler struct {
	Repo *repository.EntityRepository // Repository Injection
}

func parseAndSanitizeStats(raw json.RawMessage) (datatypes.JSON, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}

	var stats models.MonsterStats
	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&stats); err != nil {
		return nil, err
	}

	cleanBytes, err := json.Marshal(stats)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(cleanBytes), nil
}

// Create handle with create route
func (h *EntityHandler) Create(c echo.Context) error {
	var req request.CreateEntityRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Formato de dados inválido"})
	}

	name, err := request.ValidateName(req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if !req.Type.IsValid() {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tipo de local entidade inválido. Tipos permitidos: NPC, GOD, MONSTER, DEMIGOD",
		})
	}

	if err := request.ValidateLore(req.Lore); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	imageURL, err := request.ValidateImageURL(req.ImageURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	cleanStats, err := parseAndSanitizeStats(req.Stats)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "O formato dos atributos (stats) da entidade está incorreto",
		})
	}

	entity := models.Entity{
		Name:     name,
		Type:     req.Type,
		Lore:     req.Lore,
		ImageURL: imageURL,
		Stats:    cleanStats,
	}

	// Call repo to save
	if err := h.Repo.Create(&entity); err != nil {
		c.Logger().Errorf("Erro ao salvar entidade: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar entidade"})
	}

	return c.JSON(http.StatusCreated, entity)
}

// AddDrop lida com a rota POST /entities/:id/drops/:item_id
func (h *EntityHandler) AddDrop(c echo.Context) error {
	entID, err := strconv.Atoi(c.Param("id"))
	if err != nil || entID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID da entidade deve ser um número válido"})
	}

	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil || itemID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do item deve ser um número válido"})
	}

	if err := h.Repo.AddDrop(uint(entID), uint(itemID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Entidade ou item não encontrado"})
		}
		c.Logger().Errorf("Erro ao vincular drop: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao vincular drop ao monstro"})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Item adicionado aos drops do monstro com sucesso",
	})
}

// RemoveDrop lida com a rota DELETE /entities/:id/drops/:item_id
func (h *EntityHandler) RemoveDrop(c echo.Context) error {
	entID, err := strconv.Atoi(c.Param("id"))
	if err != nil || entID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID da entidade deve ser um número válido"})
	}

	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil || itemID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do item deve ser um número válido"})
	}

	if err := h.Repo.RemoveDrop(uint(entID), uint(itemID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Entidade ou item não encontrado"})
		}
		c.Logger().Errorf("Erro ao remover drop: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao remover drop da entidade"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Item removido dos drops da entidade com sucesso",
	})
}

// GetByID handle with get route /entities/:id
func (h *EntityHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número válido"})
	}

	entity, err := h.Repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Entidade não encontrada"})
		}
		c.Logger().Errorf("Erro ao buscar entidade %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro interno ao buscar entidade"})
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
	// Pega os valores da URL (ex: ?type=GOD&search=podridão)
	entityType := strings.TrimSpace(c.QueryParam("type"))
	search := strings.TrimSpace(c.QueryParam("search"))

	if len(search) > request.MaxSearchLength {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Termo de busca muito longo"})
	}

	// Passa pro repositório
	entities, err := h.Repo.GetAll(entityType, search)
	if err != nil {
		c.Logger().Errorf("Erro ao listar entidades: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao buscar entidades"})
	}

	resp := []response.MinimalResource{}
	for _, ent := range entities {
		resp = append(resp, response.MinimalResource{
			ID:   ent.ID,
			Name: ent.Name,
			Type: string(ent.Type),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// Update handle with PUT /entities/:id
func (h *EntityHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número válido"})
	}

	entity, err := h.Repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Entidade não encontrada"})
		}
		c.Logger().Errorf("Erro ao buscar entidade %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro interno ao buscar entidade"})
	}

	var req request.UpdateEntityRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Formato de dados inválido"})
	}

	if req.Name != nil {
		name, err := request.ValidateName(*req.Name)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		entity.Name = name
	}

	if req.Type != nil {
		if !req.Type.IsValid() {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Tipo de local entidade inválido. Tipos permitidos: NPC, GOD, MONSTER, DEMIGOD",
			})
		}
		entity.Type = *req.Type
	}

	if req.Lore != nil {
		if err := request.ValidateLore(*req.Lore); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		entity.Lore = *req.Lore
	}

	if req.ImageURL != nil {
		imageURL, err := request.ValidateImageURL(*req.ImageURL)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		entity.ImageURL = imageURL
	}

	if len(req.Stats) > 0 {
		cleanStats, err := parseAndSanitizeStats(req.Stats)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "O formato dos atributos (stats) da entidade está incorreto",
			})
		}
		entity.Stats = cleanStats
	}

	if err := h.Repo.Update(entity); err != nil {
		c.Logger().Errorf("Erro ao atualizar entidade %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao atualizar entidade"})
	}

	return c.JSON(http.StatusOK, entity)
}

// Delete handle with DELETE /entities/:id
func (h *EntityHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número válido"})
	}

	if err := h.Repo.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Entidade não encontrada"})
		}
		c.Logger().Errorf("Erro ao deletar entidade %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao deletar entidade"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Entidade removida com sucesso"})
}
