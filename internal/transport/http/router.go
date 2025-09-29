package http

import (
	"view-list/internal/repository"
	"view-list/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Database) *fiber.App {
	app := fiber.New()

	// Repository
	mangaRepo := repository.NewMangaRepo(db)

	// Service
	mangaSvc := service.NewMangaService(mangaRepo)

	// Handler
	mangaHandler := NewMangaHandler(mangaSvc)

	// Health route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Group de mangas
	mangaGroup := app.Group("/mangas")
	mangaGroup.Post("/", mangaHandler.CreateManga)
	mangaGroup.Get("/", mangaHandler.GetMangas)
	mangaGroup.Get("/:id", mangaHandler.GetManga)
	mangaGroup.Put("/:id", mangaHandler.UpdateManga)
	mangaGroup.Delete("/:id", mangaHandler.DeleteManga)

	return app
}
