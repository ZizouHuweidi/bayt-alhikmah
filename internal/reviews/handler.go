package reviews

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/maktaba/internal/auth"
	"github.com/zizouhuweidi/maktaba/internal/echox"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type CreateRequest struct {
	SourceID string  `json:"source_id" validate:"required"`
	Rating   int     `json:"rating" validate:"required,min=1,max=5"`
	Content  *string `json:"content,omitempty"`
	IsPublic bool    `json:"is_public"`
}

type UpdateRequest struct {
	Rating   *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Content  *string `json:"content,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
}

func (h *Handler) RegisterPublicRoutes(e *echo.Echo) {
	e.GET("/reviews", h.List)
	e.GET("/reviews/:id", h.GetByID)
}

func (h *Handler) RegisterProtectedRoutes(g *echo.Group) {
	g.POST("/reviews", h.Create)
	g.GET("/reviews", h.ListOwn)
	g.PUT("/reviews/:id", h.Update)
	g.DELETE("/reviews/:id", h.Delete)
}

func (h *Handler) Create(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	var req CreateRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}
	sourceID, err := uuid.FromString(req.SourceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid source_id")
	}

	review, err := h.service.Create(c.Request().Context(), CreateReviewParams{UserID: userID, SourceID: sourceID, Rating: req.Rating, Content: req.Content, IsPublic: req.IsPublic})
	if errors.Is(err, ErrInvalidReview) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid review")
	}
	if errors.Is(err, ErrReviewExists) {
		return echo.NewHTTPError(http.StatusConflict, "review already exists for source")
	}
	if errors.Is(err, ErrSourceNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "source not found")
	}
	if errors.Is(err, ErrReviewConflict) {
		return echo.NewHTTPError(http.StatusConflict, "review conflict")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create review")
	}

	return c.JSON(http.StatusCreated, review)
}

func (h *Handler) GetByID(c *echo.Context) error {
	id, err := echox.ParamUUID(c, "id", "review ID")
	if err != nil {
		return err
	}

	review, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrReviewNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "review not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get review")
	}
	if !review.IsPublic {
		return echo.NewHTTPError(http.StatusForbidden, "review is private")
	}

	return c.JSON(http.StatusOK, review)
}

func (h *Handler) List(c *echo.Context) error {
	limit, offset := echox.Pagination(c)
	userIDStr := c.QueryParam("user_id")
	sourceIDStr := c.QueryParam("source_id")

	var result []*Review
	var err error
	if userIDStr != "" {
		userID, parseErr := uuid.FromString(userIDStr)
		if parseErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user_id")
		}
		result, err = h.service.ListPublicByUser(c.Request().Context(), userID, limit, offset)
	} else if sourceIDStr != "" {
		sourceID, parseErr := uuid.FromString(sourceIDStr)
		if parseErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid source_id")
		}
		result, err = h.service.ListPublicBySource(c.Request().Context(), sourceID, limit, offset)
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id or source_id required")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list reviews")
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) ListOwn(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	limit, offset := echox.Pagination(c)
	reviews, err := h.service.ListByUser(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list reviews")
	}
	return c.JSON(http.StatusOK, reviews)
}

func (h *Handler) Update(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	id, err := echox.ParamUUID(c, "id", "review ID")
	if err != nil {
		return err
	}

	existing, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrReviewNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "review not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get review")
	}
	if existing.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "cannot update another user's review")
	}

	var req UpdateRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}

	review, err := h.service.Update(c.Request().Context(), id, UpdateReviewParams{Rating: req.Rating, Content: req.Content, IsPublic: req.IsPublic})
	if errors.Is(err, ErrInvalidReview) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid review")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update review")
	}
	return c.JSON(http.StatusOK, review)
}

func (h *Handler) Delete(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	id, err := echox.ParamUUID(c, "id", "review ID")
	if err != nil {
		return err
	}

	existing, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrReviewNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "review not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get review")
	}
	if existing.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "cannot delete another user's review")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete review")
	}
	return c.NoContent(http.StatusNoContent)
}
