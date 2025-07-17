package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zizouhuweidi/bayt-alhikmah/internal/repository"
)

type BookHandler struct {
	Repo repository.BookRepository
}

func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{Repo: repo}
}

func (h *BookHandler) ListBooks(c echo.Context) error {
	ctx := c.Request().Context()
	title := c.QueryParam("title")
	books, err := h.Repo.ListBooks(ctx, title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, books)
}
