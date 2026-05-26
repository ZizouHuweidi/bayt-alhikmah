package sources

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/gofrs/uuid/v5"
)

type fakeSourceRepo struct {
	createCalled bool
	existing     *Source
}

func (r *fakeSourceRepo) Create(ctx context.Context, source *Source) (*Source, error) {
	r.createCalled = true
	return source, nil
}

func (r *fakeSourceRepo) CreateBook(ctx context.Context, params CreateBookParams) (*Book, error) {
	return nil, nil
}

func (r *fakeSourceRepo) GetByID(ctx context.Context, id uuid.UUID) (*Source, error) {
	return r.existing, nil
}

func (r *fakeSourceRepo) GetBookByID(ctx context.Context, id uuid.UUID) (*Book, error) {
	return nil, nil
}

func (r *fakeSourceRepo) List(ctx context.Context, limit, offset int) ([]*Source, error) {
	return nil, nil
}

func (r *fakeSourceRepo) ListByType(ctx context.Context, sourceType SourceType, limit, offset int) ([]*Source, error) {
	return nil, nil
}

func (r *fakeSourceRepo) Update(ctx context.Context, source *Source) (*Source, error) {
	return source, nil
}

func (r *fakeSourceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *fakeSourceRepo) Search(ctx context.Context, query string, limit, offset int) ([]*Source, error) {
	return nil, nil
}

func (r *fakeSourceRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

func TestCreateRejectsInvalidSourceType(t *testing.T) {
	repo := &fakeSourceRepo{}
	service := NewService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))

	_, err := service.Create(context.Background(), CreateSourceParams{Title: "Invalid", Type: SourceType("unknown")})
	if !errors.Is(err, ErrInvalidSource) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidSource)
	}
	if repo.createCalled {
		t.Fatal("repository Create should not be called for invalid source")
	}
}

func TestUpdateRejectsInvalidSourceType(t *testing.T) {
	repo := &fakeSourceRepo{existing: &Source{ID: uuid.Must(uuid.NewV7()), Title: "Book", Type: SourceTypeBook}}
	service := NewService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))
	invalid := SourceType("unknown")

	_, err := service.Update(context.Background(), repo.existing.ID, UpdateSourceParams{Type: &invalid})
	if !errors.Is(err, ErrInvalidSource) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidSource)
	}
}
