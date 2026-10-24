package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/dto/response"
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
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.Repo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Item não encontrado"})
	}
	return c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) GetAll(c echo.Context) error {
	items, err := h.Repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao buscar itens"})
	}

	var resp []response.MinimalResource

	for _, ent := range items {
		resp = append(resp, response.MinimalResource{
			ID:   ent.ID,
			Name: ent.Name,
			Type: string(ent.Type),
		})
	}

	return c.JSON(http.StatusOK, resp)
}
