package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MangaRepo interface {
	Create(ctx context.Context, manga *Manga) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Manga, error)
	List(ctx context.Context) ([]Manga, error)
	Update(ctx context.Context, id primitive.ObjectID, updates bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
