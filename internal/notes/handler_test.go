package notes

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

	c := testContext(http.MethodGet, "/notes/"+noteID.String(), "id", noteID.String())

	err := handler.GetByID(c)

	if code := statusCode(t, err); code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, code)
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

	c := testContext(http.MethodDelete, "/api/notes/"+noteID.String(), "id", noteID.String())
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
