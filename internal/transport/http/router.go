package http

import (
	"path/filepath"
	"view-list/internal/repository"
	"view-list/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Database, staticDir string) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100 MB, si hay más tira error
	})

	// --- CORS ---
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	// --- Repository ---
	mangaRepo := repository.NewMangaRepo(db)
	userRepo := repository.NewUserRepo(db)

	// --- Services ---
	mangaSvc := service.NewMangaService(mangaRepo)
	userSvc := service.NewUserService(userRepo)

	// --- Handlers ---
	mangaHandler := NewMangaHandler(mangaSvc)
	userHandler := NewUserHandler(userSvc)

	// --- Health check ---
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// --- Auth (público) ---
	auth := app.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// --- Protected API ---
	api := app.Group("/api", JWTMiddleware())

	api.Get("/me", userHandler.Me)

	mangaGroup := api.Group("/mangas")
	mangaGroup.Post("/", mangaHandler.CreateManga)
	mangaGroup.Get("/", mangaHandler.GetMangas)
	mangaGroup.Get("/:id", mangaHandler.GetManga)
	mangaGroup.Put("/:id", mangaHandler.UpdateManga)
	mangaGroup.Delete("/:id", mangaHandler.DeleteManga)

	backupGroup := api.Group("/backup")
	backupGroup.Get("/", mangaHandler.ExportUserMangas)
	backupGroup.Post("/", mangaHandler.ImportUserMangas)

	// --- Frontend estático ---
	if staticDir != "" {
		app.Static("/", staticDir)
		app.Get("/*", func(c *fiber.Ctx) error {
			return c.SendFile(filepath.Join(staticDir, "index.html"))
		})
	}

	return app
}
