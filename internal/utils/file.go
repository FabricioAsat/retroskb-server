package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func RemoveUserUploadsAsync(userID string) {
	go func() {
		// Delay inicial para dar tiempo a liberar handles
		time.Sleep(1 * time.Second)

		dir := filepath.Join("uploads", "user_"+userID)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Printf("⚠ Error leyendo directorio %s: %v\n", dir, err)
			return
		}

		deletedCount := 0
		failedCount := 0

		// Eliminar cada archivo con retry
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			filePath := filepath.Join(dir, file.Name())

			if err := DeleteFileWithRetry(filePath, 10); err != nil {
				failedCount++
				fmt.Printf("⚠ Falló eliminar: %s\n", file.Name())
			} else {
				deletedCount++
			}
		}

		// Intentar eliminar el directorio
		time.Sleep(500 * time.Millisecond)
		for i := 0; i < 5; i++ {
			if err := os.Remove(dir); err == nil {
				fmt.Printf("✓ Directorio eliminado: %s (%d archivos)\n", dir, deletedCount)
				return
			}
			time.Sleep(time.Second * time.Duration(i+1))
		}

		fmt.Printf("⚠ Directorio no eliminado: %s (archivos: %d ok, %d fallidos)\n",
			dir, deletedCount, failedCount)
	}()
}

// Esta func xq los archivos se tardan en salir del cache
func DeleteFileWithRetry(path string, maxRetries int) error {
	if path == "" {
		return nil
	}
	p := filepath.Clean(path)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil
	}

	var lastErr error
	// Delays: 200ms, 500ms, 1s, 2s, 3s...
	delays := []time.Duration{200, 500, 1000, 2000, 3000}

	for i := 0; i < maxRetries; i++ {
		err := os.Remove(p)
		if err == nil {
			return nil
		}
		lastErr = err

		if i < len(delays) {
			time.Sleep(time.Millisecond * delays[i])
		} else {
			time.Sleep(time.Second * 3)
		}
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
