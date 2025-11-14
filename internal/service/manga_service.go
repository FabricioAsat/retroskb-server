package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"view-list/internal/domain"
	"view-list/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MangaService struct {
	mgRepo domain.MangaRepo
}

func NewMangaService(mgRepo domain.MangaRepo) *MangaService {
	return &MangaService{mgRepo: mgRepo}
}

func (s *MangaService) Create(ctx context.Context, manga *domain.Manga, userID string) error {
	// 1.0 Valido que el estado esté contemplado
	if !domain.IsValidMangaState(manga.State) {
		return errors.New("Invalid manga state")
	}

	// 1.1 Valido que el nombre no esté vacío
	if manga.Name == "" {
		return errors.New("Name cannot be empty")
	}
	// 1.2 Asigno el userID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	manga.UserID = objID

	return s.mgRepo.Create(ctx, manga)
}

func (s *MangaService) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Manga, error) {
	return s.mgRepo.GetByID(ctx, id)
}

func (s *MangaService) ListAll(ctx context.Context, userID, state, search string) ([]domain.Manga, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.mgRepo.List(ctx, objID, state, search)
}

func (s *MangaService) Update(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	// 1.0 Valido que el manga exista
	_, err := s.mgRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 1.1 Valido los states
	if val, ok := updates["state"]; ok {
		state := domain.MangaState(fmt.Sprint(val)) // Convierte a strign
		if !domain.IsValidMangaState(state) {
			return errors.New("Invalid manga state")
		}
		updates["state"] = state // Normalización
	}

	return s.mgRepo.Update(ctx, id, updates)
}

func (s *MangaService) Delete(ctx context.Context, id primitive.ObjectID) error {
	// 1.0 Valido que el manga exista
	manga, err := s.mgRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if manga.Image != "" && strings.HasPrefix(manga.Image, "/uploads/") {
		filePath := "." + manga.Image
		go func() {
			time.Sleep(200 * time.Millisecond) // Más rápido pero suficiente
			if err := utils.DeleteFileWithRetry(filePath, 8); err != nil {
				log.Printf("warning: error deleting file %s: %v\n", filePath, err)
			}
		}()
	}

	return s.mgRepo.Delete(ctx, id)
}

func (s *MangaService) DeleteAll(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// 1. borro los datos de monog
	if err := s.mgRepo.DeleteAll(ctx, objID); err != nil {
		return err
	}

	// 2. Borrar archivos de forma asíncrona (no bloquea la respuesta)
	utils.RemoveUserUploadsAsync(userID)

	return nil
}

func (s *MangaService) ExportUserMangas(ctx context.Context, userID string) ([]byte, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	mangas, err := s.mgRepo.List(ctx, objID, "", "")
	if err != nil {
		return nil, err
	}

	for i := range mangas {
		if mangas[i].Image != "" {
			b64, err := utils.ImageToBase64(mangas[i].Image)
			if err == nil {
				mangas[i].Image = b64
			}
		}
	}

	data, err := bson.Marshal(bson.M{"mangas": mangas})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *MangaService) ImportUserMangas(ctx context.Context, userID string, data []byte) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	var wrapper struct {
		Mangas []domain.Manga `bson:"mangas"`
	}

	if err := bson.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	for i := range wrapper.Mangas {
		m := &wrapper.Mangas[i]

		// ⚙️ Si viene una imagen base64, la guardamos en disco
		if strings.HasPrefix(m.Image, "data:image/") {
			// 🧹 Limpiar posibles saltos de línea o espacios
			clean := strings.ReplaceAll(m.Image, "\n", "")
			clean = strings.ReplaceAll(clean, "\r", "")
			clean = strings.TrimSpace(clean)

			imgPath, err := utils.SaveBase64ImageForUser(clean, userID)
			if err != nil {
				fmt.Printf("❌ Error saving image for manga %s: %v\n", m.Name, err)
				m.Image = "" // limpiar si falló
			} else {
				m.Image = imgPath
			}
		}

		// limpiar IDs y asignar usuario actual
		m.ID = primitive.NilObjectID
		m.UserID = objID
	}

	// insertar todos
	if err := s.mgRepo.BulkInsert(ctx, wrapper.Mangas); err != nil {
		return err
	}

	return nil
}
