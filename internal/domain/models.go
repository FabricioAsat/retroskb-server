package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MangaState string

const (
	MangaStateReading   MangaState = "reading"
	MangaStateCompleted MangaState = "completed"
	MangaStateAbandoned MangaState = "abandoned"
	MangaStateDeleted   MangaState = "deleted"
	MangaStateOnHold    MangaState = "on hold"
)

type Manga struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name      string             `bson:"name,omitempty" json:"name"`
	State     MangaState         `bson:"state" json:"state"`
	Chapter   uint16             `bson:"chapter,min=0" json:"chapter"` // cap en el que lo dejé
	Image     []byte             `bson:"image" json:"image"`
	Link      string             `bson:"link" json:"link"` // link donde lo miro
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func IsValidMangaState(s MangaState) bool {
	switch s {
	case MangaStateReading,
		MangaStateCompleted,
		MangaStateAbandoned,
		MangaStateDeleted,
		MangaStateOnHold:
		return true
	}
	return false
}
