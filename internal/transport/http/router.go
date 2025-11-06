package http

import (
	"view-list/internal/repository"
	"view-list/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Database) *fiber.App {
	app := fiber.New()

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

	// --- Public routes ---
	auth := app.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// --- Protected routes ---
	protected := app.Group("/", JWTMiddleware())
	protected.Get("/me", userHandler.Me)

	// --- Mangas (protected) ---
	mangaGroup := protected.Group("/mangas")
	mangaGroup.Post("/", mangaHandler.CreateManga)
	mangaGroup.Get("/", mangaHandler.GetMangas)
	mangaGroup.Get("/:id", mangaHandler.GetManga)
	mangaGroup.Put("/:id", mangaHandler.UpdateManga)
	mangaGroup.Delete("/:id", mangaHandler.DeleteManga)

	return app
}
