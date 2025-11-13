package service

import (
	"context"
	"errors"
	"fmt"
	"view-list/internal/domain"

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

func (s *MangaService) ListAll(ctx context.Context, userID string) ([]domain.Manga, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.mgRepo.List(ctx, objID)
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
	_, err := s.mgRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.mgRepo.Delete(ctx, id)
}

func (s *MangaService) ExportUserMangas(ctx context.Context, userID string) ([]byte, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	mangas, err := s.mgRepo.List(ctx, objID)
	if err != nil {
		return nil, err
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
		wrapper.Mangas[i].ID = primitive.NilObjectID
		wrapper.Mangas[i].UserID = objID
	}

	if err := s.mgRepo.BulkInsert(ctx, wrapper.Mangas); err != nil {
		return err
	}

	return nil
}
