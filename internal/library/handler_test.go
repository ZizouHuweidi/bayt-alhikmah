package library

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
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

	req := httptest.NewRequest(http.MethodGet, "/api/library/items/"+itemID.String(), nil)
	req = req.WithContext(auth.ContextWithUserID(req.Context(), requesterID))
	req.SetPathValue("id", itemID.String())
	res := httptest.NewRecorder()

	handler.GetMine(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
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

	req := httptest.NewRequest(http.MethodDelete, "/api/library/items/"+itemID.String(), nil)
	req = req.WithContext(auth.ContextWithUserID(req.Context(), requesterID))
	req.SetPathValue("id", itemID.String())
	res := httptest.NewRecorder()

	handler.Delete(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
	}
	if repo.deleted {
		t.Fatal("expected delete not to be called")
	}
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
