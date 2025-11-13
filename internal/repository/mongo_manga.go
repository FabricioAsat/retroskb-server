package repository

import (
	"context"
	"view-list/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoMangaRepo struct {
	db *mongo.Collection
}

func NewMangaRepo(db *mongo.Database) domain.MangaRepo {
	return &MongoMangaRepo{db: db.Collection("mangas")}
}

func (r *MongoMangaRepo) Create(ctx context.Context, manga *domain.Manga) error {
	_, err := r.db.InsertOne(ctx, manga)
	return err
}

func (r *MongoMangaRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Manga, error) {
	var manga domain.Manga
	if err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&manga); err != nil {
		return nil, err
	}
	return &manga, nil
}

// Esto trae por user_id mediante jwt, no me trae todos,
func (r *MongoMangaRepo) List(ctx context.Context, userID primitive.ObjectID) ([]domain.Manga, error) {
	filter := bson.M{"user_id": userID}
	var mangas []domain.Manga
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &mangas); err != nil {
		return nil, err
	}
	return mangas, nil
}

func (r *MongoMangaRepo) Update(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

func (r *MongoMangaRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Este inserta todos los mangas del binario
func (r *MongoMangaRepo) BulkInsert(ctx context.Context, mangas []domain.Manga) error {
	docs := make([]any, len(mangas))
	for i, m := range mangas {
		docs[i] = m
	}
	_, err := r.db.InsertMany(ctx, docs)
	return err
}
