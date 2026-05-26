package collections

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
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

	req := httptest.NewRequest(http.MethodGet, "/collections/"+collectionID.String(), nil)
	req.SetPathValue("id", collectionID.String())
	res := httptest.NewRecorder()

	handler.GetByID(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
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

	req := httptest.NewRequest(http.MethodDelete, "/api/collections/"+collectionID.String(), nil)
	req = req.WithContext(auth.ContextWithUserID(req.Context(), mustTestUUID(t)))
	req.SetPathValue("id", collectionID.String())
	res := httptest.NewRecorder()

	handler.Delete(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
	}
	if repo.deleted {
		t.Fatal("expected delete not to be called")
	}
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
