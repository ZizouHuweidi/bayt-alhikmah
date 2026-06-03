package library

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/maktaba/internal/auth"
)

func TestHandlerGetMineRejectsAnotherUsersItem(t *testing.T) {
	itemID := mustTestUUID(t)
	ownerID := mustTestUUID(t)
	requesterID := mustTestUUID(t)
	handler := NewHandler(NewService(&fakeLibraryRepository{item: &Item{
		ID:         itemID,
		UserID:     ownerID,
		SourceID:   mustTestUUID(t),
		Status:     StatusToConsume,
		Visibility: VisibilityPrivate,
	}}, slog.Default()), slog.Default())

	c := testContext(http.MethodGet, "/api/library/items/"+itemID.String(), "id", itemID.String())
	auth.SetUserID(c, requesterID)

	err := handler.GetMine(c)

	if code := statusCode(t, err); code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, code)
	}
}

func TestHandlerDeleteRejectsAnotherUsersItem(t *testing.T) {
	itemID := mustTestUUID(t)
	ownerID := mustTestUUID(t)
	requesterID := mustTestUUID(t)
	repo := &fakeLibraryRepository{item: &Item{
		ID:         itemID,
		UserID:     ownerID,
		SourceID:   mustTestUUID(t),
		Status:     StatusToConsume,
		Visibility: VisibilityPrivate,
	}}
	handler := NewHandler(NewService(repo, slog.Default()), slog.Default())

	c := testContext(http.MethodDelete, "/api/library/items/"+itemID.String(), "id", itemID.String())
	auth.SetUserID(c, requesterID)

	err := handler.Delete(c)

	if code := statusCode(t, err); code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, code)
	}
	if repo.deleted {
		t.Fatal("expected delete not to be called")
	}
}

func testContext(method, target, paramName, paramValue string) *echo.Context {
	e := echo.New()
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{{Name: paramName, Value: paramValue}})
	return c
}

func statusCode(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}
	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected echo HTTPError, got %T", err)
	}
	return httpErr.Code
}

type fakeLibraryRepository struct {
	item    *Item
	deleted bool
}

func (r *fakeLibraryRepository) Create(context.Context, *Item) (*Item, error) {
	panic("not implemented")
}

func (r *fakeLibraryRepository) GetByID(_ context.Context, id uuid.UUID) (*Item, error) {
	if r.item == nil || r.item.ID != id {
		return nil, nil
	}
	return r.item, nil
}

func (r *fakeLibraryRepository) ListByUser(context.Context, uuid.UUID, int, int) ([]*Item, error) {
	panic("not implemented")
}

func (r *fakeLibraryRepository) ListPublicByUser(context.Context, uuid.UUID, int, int) ([]*Item, error) {
	panic("not implemented")
}

func (r *fakeLibraryRepository) ListPublicByUsername(context.Context, string, int, int) ([]*Item, error) {
	panic("not implemented")
}

func (r *fakeLibraryRepository) ListByUserWithSources(context.Context, uuid.UUID, int, int) ([]*ItemWithSource, error) {
	panic("not implemented")
}

func (r *fakeLibraryRepository) ListPublicByUsernameWithSources(context.Context, string, int, int) ([]*ItemWithSource, error) {
	panic("not implemented")
}

func (r *fakeLibraryRepository) Update(context.Context, *Item) (*Item, error) {
	panic("not implemented")
}

func (r *fakeLibraryRepository) Delete(context.Context, uuid.UUID) error {
	r.deleted = true
	return nil
}

func mustTestUUID(t *testing.T) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("new uuid: %v", err)
	}
	return id
}
