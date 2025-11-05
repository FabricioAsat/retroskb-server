package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MangaRepo interface {
	Create(ctx context.Context, manga *Manga) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Manga, error)
	List(ctx context.Context, userID primitive.ObjectID) ([]Manga, error)
	Update(ctx context.Context, id primitive.ObjectID, updates bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// -------------------- USERS --------------------

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type UserService interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, email, password string) (*User, error)
}
