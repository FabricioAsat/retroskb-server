package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"view-list/internal/domain"
	"view-list/internal/utils" // asegurate que exista utils.SaveBase64Image
)

// Funcion usada para migrar de base64 a url
func main() {
	ctx := context.Background()

	// ⚙️ Conectarse a Mongo
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("retroskb")
	coll := db.Collection("mangas")

	// 🔍 Buscar todos los mangas
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var mangas []domain.Manga
	if err := cursor.All(ctx, &mangas); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Encontrados %d mangas\n", len(mangas))

	// 🔁 Recorrer y migrar
	for _, m := range mangas {
		if strings.HasPrefix(m.Image, "data:image") {
			url, err := utils.SaveBase64Image(m.Image)
			if err != nil {
				log.Printf("❌ Error guardando imagen de %s: %v\n", m.Name, err)
				continue
			}

			_, err = coll.UpdateOne(ctx,
				bson.M{"_id": m.ID},
				bson.M{"$set": bson.M{"image": url}},
			)
			if err != nil {
				log.Printf("⚠️ Error actualizando %s: %v\n", m.Name, err)
				continue
			}

			log.Printf("✅ Migrado: %s → %s\n", m.Name, url)
		}
	}

	fmt.Println("🎉 Migración completa.")
}
