package reviews

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

func TestHandlerGetByIDRejectsPrivateReview(t *testing.T) {
	reviewID := mustTestUUID(t)
	handler := NewHandler(NewService(&fakeReviewsRepository{review: &Review{
		ID:       reviewID,
		UserID:   mustTestUUID(t),
		SourceID: mustTestUUID(t),
		Rating:   4,
		IsPublic: false,
	}}, slog.Default()), slog.Default())

	c := testContext(http.MethodGet, "/reviews/"+reviewID.String(), "id", reviewID.String())

	err := handler.GetByID(c)

	if code := statusCode(t, err); code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, code)
	}
}

func TestHandlerDeleteRejectsAnotherUsersReview(t *testing.T) {
	reviewID := mustTestUUID(t)
	repo := &fakeReviewsRepository{review: &Review{
		ID:       reviewID,
		UserID:   mustTestUUID(t),
		SourceID: mustTestUUID(t),
		Rating:   4,
		IsPublic: true,
	}}
	handler := NewHandler(NewService(repo, slog.Default()), slog.Default())

	c := testContext(http.MethodDelete, "/api/reviews/"+reviewID.String(), "id", reviewID.String())
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

type fakeReviewsRepository struct {
	review  *Review
	deleted bool
}

func (r *fakeReviewsRepository) Create(context.Context, *Review) (*Review, error) {
	panic("not implemented")
}

func (r *fakeReviewsRepository) GetByID(_ context.Context, id uuid.UUID) (*Review, error) {
	if r.review == nil || r.review.ID != id {
		return nil, nil
	}
	return r.review, nil
}

func (r *fakeReviewsRepository) ListByUser(context.Context, uuid.UUID, int, int) ([]*Review, error) {
	panic("not implemented")
}

func (r *fakeReviewsRepository) ListBySource(context.Context, uuid.UUID, int, int) ([]*Review, error) {
	panic("not implemented")
}

func (r *fakeReviewsRepository) ListPublicByUser(context.Context, uuid.UUID, int, int) ([]*Review, error) {
	panic("not implemented")
}

func (r *fakeReviewsRepository) ListPublicBySource(context.Context, uuid.UUID, int, int) ([]*Review, error) {
	panic("not implemented")
}

func (r *fakeReviewsRepository) Update(context.Context, *Review) (*Review, error) {
	panic("not implemented")
}

func (r *fakeReviewsRepository) Delete(context.Context, uuid.UUID) error {
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
