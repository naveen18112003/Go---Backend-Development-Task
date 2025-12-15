package handler

import (
	"net/http"
	"strconv"

	"user-age-api/internal/models"
	"user-age-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserHandler wires HTTP requests to the service layer.
type UserHandler struct {
	service *service.UserService
	logger  *zap.Logger
}

// NewUserHandler builds a handler.
func NewUserHandler(service *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

// Register attaches routes to the router group.
func (h *UserHandler) Register(r fiber.Router) {
	r.Post("/users", h.CreateUser)
	r.Get("/users", h.ListUsers)
	r.Get("/users/:id", h.GetUser)
	r.Put("/users/:id", h.UpdateUser)
	r.Delete("/users/:id", h.DeleteUser)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid body", zap.Error(err))
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		h.logger.Warn("failed to create user", zap.Error(err))
		return handleError(err)
	}
	return c.Status(http.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	user, err := h.service.GetUser(c.Context(), id)
	if err != nil {
		h.logger.Warn("failed to fetch user", zap.Error(err))
		return handleError(err)
	}
	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}

	user, err := h.service.UpdateUser(c.Context(), id, req)
	if err != nil {
		h.logger.Warn("failed to update user", zap.Error(err))
		return handleError(err)
	}
	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	if err := h.service.DeleteUser(c.Context(), id); err != nil {
		h.logger.Warn("failed to delete user", zap.Error(err))
		return handleError(err)
	}
	return c.SendStatus(http.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit := int32(50)
	offset := int32(0)

	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = int32(v)
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = int32(v)
		}
	}

	users, err := h.service.ListUsers(c.Context(), limit, offset)
	if err != nil {
		h.logger.Warn("failed to list users", zap.Error(err))
		return handleError(err)
	}
	return c.JSON(users)
}

func parseID(idStr string) (int32, error) {
	val, err := strconv.ParseInt(idStr, 10, 32)
	return int32(val), err
}

func handleError(err error) error {
	switch err {
	case service.ErrNotFound:
		return fiber.NewError(http.StatusNotFound, err.Error())
	default:
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
}



