package handlers

import (
	"net/http"
	"strconv"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/repository"
	"github.com/labstack/echo/v4"
)

type ItemHandler struct {
	Repo *repository.ItemRepository
}

func (h *ItemHandler) Create(c echo.Context) error {
	var item models.Item
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dados inválidos"})
	}
	if err := h.Repo.Create(&item); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar"})
	}
	return c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.Repo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Item não encontrado"})
	}
	return c.JSON(http.StatusOK, item)
}
