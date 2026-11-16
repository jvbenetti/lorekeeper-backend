package request

import (
	"errors"
	"net/url"
	"strings"
)

const (
	MaxNameLength     = 255
	MaxImageURLLength = 500
	MaxLoreLength     = 50000
	MaxSearchLength   = 100
)

func ValidateName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", errors.New("O campo 'name' é obrigatório")
	}
	if len(trimmed) > MaxNameLength {
		return "", errors.New("O campo 'name' deve ter no máximo 255 caracteres")
	}
	return trimmed, nil
}

func ValidateLore(lore string) error {
	if len(lore) > MaxLoreLength {
		return errors.New("O campo 'lore' excede o tamanho máximo permitido")
	}
	return nil
}

func ValidateImageURL(rawURL string) (string, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", nil
	}
	if len(trimmed) > MaxImageURLLength {
		return "", errors.New("O campo 'image_url' deve ter no máximo 500 caracteres")
	}
	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", errors.New("O campo 'image_url' deve ser uma URL válida (http ou https)")
	}
	return trimmed, nil
}
