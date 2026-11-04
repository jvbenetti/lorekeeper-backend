package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login valida as credenciais estáticas e retorna um JWT
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Requisição inválida"})
	}

	expectedUser := os.Getenv("ADMIN_USERNAME")
	expectedPass := os.Getenv("ADMIN_PASSWORD")
	jwtSecret := os.Getenv("JWT_SECRET")

	if expectedUser == "" || expectedPass == "" || jwtSecret == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Credenciais ou JWT_SECRET não configurados no servidor"})
	}

	// Valida usuário e senha
	if req.Username != expectedUser || req.Password != expectedPass {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Credenciais incorretas"})
	}

	// Gera o token válido por 72 horas
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		Subject:   req.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao gerar token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
