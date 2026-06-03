package collections

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

func TestHandlerGetByIDRejectsPrivateCollection(t *testing.T) {
	collectionID := mustTestUUID(t)
	handler := NewHandler(NewService(&fakeCollectionsRepository{collection: &Collection{
		ID:       collectionID,
		UserID:   mustTestUUID(t),
		Name:     "Private list",
		IsPublic: false,
	}}, slog.Default()), slog.Default())

	c := testContext(http.MethodGet, "/collections/"+collectionID.String(), "id", collectionID.String())

	err := handler.GetByID(c)

	if code := statusCode(t, err); code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, code)
	}
}

func TestHandlerDeleteRejectsAnotherUsersCollection(t *testing.T) {
	collectionID := mustTestUUID(t)
	repo := &fakeCollectionsRepository{collection: &Collection{
		ID:       collectionID,
		UserID:   mustTestUUID(t),
		Name:     "Someone else's list",
		IsPublic: true,
	}}
	handler := NewHandler(NewService(repo, slog.Default()), slog.Default())

	c := testContext(http.MethodDelete, "/api/collections/"+collectionID.String(), "id", collectionID.String())
	auth.SetUserID(c, mustTestUUID(t))

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

type fakeCollectionsRepository struct {
	collection *Collection
	deleted    bool
}

func (r *fakeCollectionsRepository) Create(context.Context, *Collection) (*Collection, error) {
	panic("not implemented")
}

func (r *fakeCollectionsRepository) GetByID(_ context.Context, id uuid.UUID) (*Collection, error) {
	if r.collection == nil || r.collection.ID != id {
		return nil, nil
	}
	return r.collection, nil
}

func (r *fakeCollectionsRepository) ListByUser(context.Context, uuid.UUID, int, int) ([]*Collection, error) {
	panic("not implemented")
}

func (r *fakeCollectionsRepository) ListPublicByUser(context.Context, uuid.UUID, int, int) ([]*Collection, error) {
	panic("not implemented")
}

func (r *fakeCollectionsRepository) Update(context.Context, *Collection) (*Collection, error) {
	panic("not implemented")
}

func (r *fakeCollectionsRepository) Delete(context.Context, uuid.UUID) error {
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
