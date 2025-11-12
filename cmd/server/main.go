package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"view-list/internal/transport/http"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, url)
	exec.Command(cmd, args...).Start()
}
func main() {
	staticDir := ""
	// 1. Cargar .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "prod"
	}

	if env == "dev" {
		log.Println("Running in development mode")
	} else {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		baseDir := filepath.Dir(exePath)
		relativePath := filepath.Join(".", "web", "dist")
		if _, err := os.Stat(relativePath); err == nil {
			staticDir = relativePath
		} else {
			staticDir = filepath.Join(baseDir, "web", "dist")
		}
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "4090"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "retroskb"
	}

	// 2. Conectar Mongo
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	db := client.Database(dbName)

	// 3. Crear router principal
	app := http.NewRouter(db, staticDir)

	// 4. Iniciar servidor y abrir navegador
	if env == "prod" {
		url := "http://localhost:" + port
		log.Println("Servidor iniciado en:", url)
		go openBrowser(url)
	}

	log.Fatal(app.Listen(":" + port))
}
