package main

import (
	"context"
	"log"
	"os"
	"view-list/internal/transport/http"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Cargar .env al entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4096"
	}

	// Conectar Mongo
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("manga-viwer")

	// Creo router
	app := http.NewRouter(db)

	// Iniciar server
	log.Println("Servidor iniciado en el puerto: ", port)
	log.Fatal(app.Listen(":" + port))

}
