package http

import (
	"view-list/internal/repository"
	"view-list/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Database) *fiber.App {
	app := fiber.New()

	// CORS Policy
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}
		return c.Next()
	})

	// Repository
	mangaRepo := repository.NewMangaRepo(db)
	userRepo := repository.NewUserRepo(db)

	// Service
	mangaSvc := service.NewMangaService(mangaRepo)
	userSvc := service.NewUserService(userRepo)

	// Handler
	mangaHandler := NewMangaHandler(mangaSvc)
	userHandler := NewUserHandler(userSvc)

	// Health route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Rutas del user
	auth := app.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// Rutas protegidas
	protected := app.Group("/", JWTMiddleware())
	protected.Get("/me", userHandler.Me)

	// Group de mangas
	mangaGroup := app.Group("/mangas")
	mangaGroup.Post("/", mangaHandler.CreateManga)
	mangaGroup.Get("/", mangaHandler.GetMangas)
	mangaGroup.Get("/:id", mangaHandler.GetManga)
	mangaGroup.Put("/:id", mangaHandler.UpdateManga)
	mangaGroup.Delete("/:id", mangaHandler.DeleteManga)

	return app
}
