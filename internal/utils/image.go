package utils

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// guarda una imagen base64 dentro de la carpeta del usuario y devuelve la URL pública.
func SaveBase64ImageForUser(base64Data, userID string) (string, error) {
	backendURL := os.Getenv("BACKEND_URL_WITHOUT_PORT") + os.Getenv("PORT")
	if base64Data == "" {
		return "", nil
	}

	// limpiar espacios, saltos o comillas que pueden romper el formato
	base64Data = strings.TrimSpace(base64Data)
	base64Data = strings.Trim(base64Data, "\"")

	// debe tener una coma que separe el header del contenido
	if !strings.Contains(base64Data, ",") {
		return "", fmt.Errorf("invalid base64 image format")
	}

	parts := strings.SplitN(base64Data, ",", 2)
	header := parts[0]
	data := parts[1]

	var ext string
	switch {
	case strings.Contains(header, "image/jpeg"):
		ext = ".jpg"
	case strings.Contains(header, "image/png"):
		ext = ".png"
	case strings.Contains(header, "image/webp"):
		ext = ".webp"
	default:
		ext = ".jpg"
	}

	// decodificar base64
	imgData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %w", err)
	}

	// crear carpeta del usuario si no existe
	dir := filepath.Join("uploads", "user_"+userID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// crear archivo con nombre único
	filename := uuid.New().String() + ext
	fullPath := filepath.Join(dir, filename)

	if err := os.WriteFile(fullPath, imgData, 0644); err != nil {
		return "", err
	}

	return backendURL + "/" + filepath.ToSlash(fullPath), nil
}

func ImageToBase64(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	// Extraer el path de la URL
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	// Remover prefijo /uploads/ si tu carpeta local se llama "uploads"
	localPath := strings.TrimPrefix(parsedURL.Path, "/")

	// Limpiar y leer el archivo
	p := filepath.Clean(localPath)
	data, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", p, err)
	}

	ext := strings.ToLower(filepath.Ext(p))

	var mime string
	switch ext {
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".png":
		mime = "image/png"
	case ".webp":
		mime = "image/webp"
	default:
		mime = "application/octet-stream"
	}

	base64Str := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mime, base64Str), nil
}
