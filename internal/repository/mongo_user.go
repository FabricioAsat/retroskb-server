package repository

import (
	"context"
	"view-list/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) domain.UserRepo {
	return &MongoUserRepo{collection: db.Collection("users")}
}

func (r *MongoUserRepo) Create(ctx context.Context, user *domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
