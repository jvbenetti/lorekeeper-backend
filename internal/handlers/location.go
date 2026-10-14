package handlers

import (
	"net/http"
	"strconv"

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

func (h *LocationHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	location, err := h.Repo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Local não encontrado"})
	}
	return c.JSON(http.StatusOK, location)
}
