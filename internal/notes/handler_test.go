package notes

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/zizouhuweidi/maktaba/internal/auth"
)

func TestHandlerGetByIDRejectsPrivateNote(t *testing.T) {
	noteID := mustTestUUID(t)
	ownerID := mustTestUUID(t)
	handler := NewHandler(NewService(&fakeNotesRepository{note: &Note{
		ID:          noteID,
		UserID:      ownerID,
		Content:     "private",
		ContentType: ContentTypeNote,
		IsPublic:    false,
	}}, slog.Default()), slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/notes/"+noteID.String(), nil)
	req.SetPathValue("id", noteID.String())
	res := httptest.NewRecorder()

	handler.GetByID(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
	}
}

func TestHandlerDeleteRejectsAnotherUsersNote(t *testing.T) {
	noteID := mustTestUUID(t)
	ownerID := mustTestUUID(t)
	requesterID := mustTestUUID(t)
	repo := &fakeNotesRepository{note: &Note{
		ID:          noteID,
		UserID:      ownerID,
		Content:     "owned by someone else",
		ContentType: ContentTypeNote,
		IsPublic:    true,
	}}
	handler := NewHandler(NewService(repo, slog.Default()), slog.Default())

	req := httptest.NewRequest(http.MethodDelete, "/api/notes/"+noteID.String(), nil)
	req = req.WithContext(auth.ContextWithUserID(req.Context(), requesterID))
	req.SetPathValue("id", noteID.String())
	res := httptest.NewRecorder()

	handler.Delete(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
	}
	if repo.deleted {
		t.Fatal("expected delete not to be called")
	}
}

type fakeNotesRepository struct {
	note    *Note
	deleted bool
}

func (r *fakeNotesRepository) Create(context.Context, *Note) (*Note, error) {
	panic("not implemented")
}

func (r *fakeNotesRepository) GetByID(_ context.Context, id uuid.UUID) (*Note, error) {
	if r.note == nil || r.note.ID != id {
		return nil, nil
	}
	return r.note, nil
}

func (r *fakeNotesRepository) ListByUser(context.Context, uuid.UUID, int, int) ([]*Note, error) {
	panic("not implemented")
}

func (r *fakeNotesRepository) ListBySource(context.Context, uuid.UUID, int, int) ([]*Note, error) {
	panic("not implemented")
}

func (r *fakeNotesRepository) ListPublicByUser(context.Context, uuid.UUID, int, int) ([]*Note, error) {
	panic("not implemented")
}

func (r *fakeNotesRepository) ListPublicBySource(context.Context, uuid.UUID, int, int) ([]*Note, error) {
	panic("not implemented")
}

func (r *fakeNotesRepository) ListPublic(context.Context, int, int) ([]*Note, error) {
	panic("not implemented")
}

func (r *fakeNotesRepository) Update(context.Context, *Note) (*Note, error) {
	panic("not implemented")
}

func (r *fakeNotesRepository) Delete(context.Context, uuid.UUID) error {
	r.deleted = true
	return nil
}

func (r *fakeNotesRepository) CountByUser(context.Context, uuid.UUID) (int64, error) {
	panic("not implemented")
}

func mustTestUUID(t *testing.T) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("new uuid: %v", err)
	}
	return id
}
