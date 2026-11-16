package handlers

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Usuário e senha são obrigatórios"})
	}

	expectedUser := os.Getenv("ADMIN_USERNAME")
	expectedPass := os.Getenv("ADMIN_PASSWORD")
	expectedPassHash := os.Getenv("ADMIN_PASSWORD_HASH")
	jwtSecret := os.Getenv("JWT_SECRET")

	if expectedUser == "" || (expectedPass == "" && expectedPassHash == "") || jwtSecret == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Credenciais ou JWT_SECRET não configurados no servidor"})
	}

	// Comparação em tempo constante para evitar Timing Attacks
	userMatch := subtle.ConstantTimeCompare([]byte(req.Username), []byte(expectedUser)) == 1

	var passMatch bool
	if expectedPassHash != "" {
		passMatch = bcrypt.CompareHashAndPassword([]byte(expectedPassHash), []byte(req.Password)) == nil
	} else {
		passMatch = subtle.ConstantTimeCompare([]byte(req.Password), []byte(expectedPass)) == 1
	}

	if !userMatch || !passMatch {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Credenciais incorretas"})
	}

	now := time.Now()
	claims := &jwt.RegisteredClaims{
		Issuer:    "lorekeeper-api",
		Subject:   expectedUser,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 12)),
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
