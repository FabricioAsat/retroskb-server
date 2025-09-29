package http

import (
	"time"
	"view-list/internal/domain"
	"view-list/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MangaHandler struct {
	svc *service.MangaService
}

func NewMangaHandler(svc *service.MangaService) *MangaHandler {
	return &MangaHandler{svc}
}

// Structs para create & update
type CreateMangaRequest struct {
	Name    string            `json:"name"`
	State   domain.MangaState `json:"state"`
	Chapter uint16            `json:"chapter"`
	Image   []byte            `json:"image"` // Mongo deja hasta 16MB por data
	Link    string            `json:"link"`
}
type UpdateMangaRequest struct {
	Name    *string            `json:"name,omitempty"`
	State   *domain.MangaState `json:"state,omitempty"`
	Chapter *uint16            `json:"chapter,omitempty"`
	Image   *[]byte            `json:"image,omitempty"`
	Link    *string            `json:"link,omitempty"`
}

// * Comienzan los métodos de la API
func (h *MangaHandler) CreateManga(c *fiber.Ctx) error {
	var req CreateMangaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}

	manga := &domain.Manga{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		State:     req.State,
		Chapter:   req.Chapter,
		Image:     req.Image,
		Link:      req.Link,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := h.svc.Create(c.Context(), manga); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": manga, "message": "Manga created successfully!"})
}

func (h *MangaHandler) GetMangas(c *fiber.Ctx) error {
	mangas, err := h.svc.ListAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": mangas, "message": "Mangas retrieved successfully!"})
}

func (h *MangaHandler) GetManga(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	manga, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": manga, "message": "Manga retrieved successfully!"})
}

func (h *MangaHandler) UpdateManga(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var req UpdateMangaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Hago el mapeo de updates
	updates := bson.M{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.State != nil {
		updates["state"] = *req.State
	}
	if req.Chapter != nil {
		updates["chapter"] = *req.Chapter
	}
	if req.Image != nil {
		updates["image"] = *req.Image
	}
	if req.Link != nil {
		updates["link"] = *req.Link
	}
	updates["updated_at"] = time.Now()

	if err := h.svc.Update(c.Context(), id, updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Manga updated successfully!"})
}

func (h *MangaHandler) DeleteManga(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Manga deleted successfully!"})
}
