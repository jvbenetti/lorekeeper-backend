package handlers

import (
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

	// Call repo to save
	if err := h.Repo.Create(&entity); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar entidade"})
	}

	return c.JSON(http.StatusCreated, entity)
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
