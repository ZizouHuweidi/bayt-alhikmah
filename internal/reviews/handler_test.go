package reviews

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
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

	req := httptest.NewRequest(http.MethodGet, "/reviews/"+reviewID.String(), nil)
	req.SetPathValue("id", reviewID.String())
	res := httptest.NewRecorder()

	handler.GetByID(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
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

	req := httptest.NewRequest(http.MethodDelete, "/api/reviews/"+reviewID.String(), nil)
	req = req.WithContext(auth.ContextWithUserID(req.Context(), mustTestUUID(t)))
	req.SetPathValue("id", reviewID.String())
	res := httptest.NewRecorder()

	handler.Delete(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
	}
	if repo.deleted {
		t.Fatal("expected delete not to be called")
	}
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
