package library

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/gofrs/uuid/v5"
)

type fakeLibraryRepo struct {
	created    *Item
	listLimit  int
	listOffset int
}

func (r *fakeLibraryRepo) Create(ctx context.Context, item *Item) (*Item, error) {
	r.created = item
	item.ID = uuid.Must(uuid.NewV7())
	return item, nil
}

func (r *fakeLibraryRepo) GetByID(ctx context.Context, id uuid.UUID) (*Item, error) {
	return nil, nil
}

func (r *fakeLibraryRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error) {
	r.listLimit = limit
	r.listOffset = offset
	return nil, nil
}

func (r *fakeLibraryRepo) ListByUserWithSources(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ItemWithSource, error) {
	r.listLimit = limit
	r.listOffset = offset
	return nil, nil
}

func (r *fakeLibraryRepo) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error) {
	return nil, nil
}

func (r *fakeLibraryRepo) ListPublicByUsername(ctx context.Context, username string, limit, offset int) ([]*Item, error) {
	return nil, nil
}

func (r *fakeLibraryRepo) ListPublicByUsernameWithSources(ctx context.Context, username string, limit, offset int) ([]*ItemWithSource, error) {
	return nil, nil
}

func (r *fakeLibraryRepo) Update(ctx context.Context, item *Item) (*Item, error) {
	return item, nil
}

func (r *fakeLibraryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestCreateDefaultsVisibilityToPrivate(t *testing.T) {
	repo := &fakeLibraryRepo{}
	service := NewService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))

	created, err := service.Create(context.Background(), CreateItemParams{
		UserID:   uuid.Must(uuid.NewV7()),
		SourceID: uuid.Must(uuid.NewV7()),
		Status:   StatusToConsume,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.Visibility != VisibilityPrivate {
		t.Fatalf("Visibility = %q, want %q", created.Visibility, VisibilityPrivate)
	}
	if repo.created == nil {
		t.Fatal("repository Create was not called")
	}
}

func TestCreateRejectsInvalidProgress(t *testing.T) {
	service := NewService(&fakeLibraryRepo{}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	negative := -1

	_, err := service.Create(context.Background(), CreateItemParams{
		UserID:        uuid.Must(uuid.NewV7()),
		SourceID:      uuid.Must(uuid.NewV7()),
		Status:        StatusInProgress,
		ProgressValue: &negative,
	})
	if !errors.Is(err, ErrInvalidItem) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidItem)
	}
}

func TestListByUserNormalizesPagination(t *testing.T) {
	repo := &fakeLibraryRepo{}
	service := NewService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))

	_, err := service.ListByUser(context.Background(), uuid.Must(uuid.NewV7()), 500, -20)
	if err != nil {
		t.Fatalf("ListByUser returned error: %v", err)
	}
	if repo.listLimit != 100 {
		t.Fatalf("limit = %d, want 100", repo.listLimit)
	}
	if repo.listOffset != 0 {
		t.Fatalf("offset = %d, want 0", repo.listOffset)
	}
}
