package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/repository"
	"github.com/labstack/echo/v4"
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

// GetByID handle with get route /entities/:id
func (h *EntityHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID deve ser um número"})
	}

	entity, err := h.Repo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Entidade não encontrada"})
	}

	return c.JSON(http.StatusOK, entity)
}
