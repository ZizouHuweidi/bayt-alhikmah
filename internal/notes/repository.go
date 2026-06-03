package notes

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zizouhuweidi/maktaba/internal/db"
	"github.com/zizouhuweidi/maktaba/internal/db/dbgen"
)

type postgresRepository struct {
	queries *dbgen.Queries
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{queries: dbgen.New(d.Pool)}
}

func (r *postgresRepository) Create(ctx context.Context, n *Note) (*Note, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	annotations, err := json.Marshal(n.Annotations)
	if err != nil {
		return nil, err
	}
	tags, err := json.Marshal(n.Tags)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.CreateNote(ctx, dbgen.CreateNoteParams{
		ID:          db.PGUUID(id),
		UserID:      db.PGUUID(n.UserID),
		SourceID:    pgUUIDPtr(n.SourceID),
		Content:     n.Content,
		ContentType: string(n.ContentType),
		IsPublic:    db.PGBool(n.IsPublic),
		Annotations: annotations,
		Tags:        tags,
	})
	if err != nil {
		return nil, mapCreateError(err)
	}
	return mapNote(row), nil
}

func mapCreateError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	if pgErr.Code == "23503" && pgErr.ConstraintName == "notes_source_id_fkey" {
		return ErrSourceNotFound
	}
	return err
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Note, error) {
	row, err := r.queries.GetNoteByID(ctx, db.PGUUID(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapNote(row), nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Note, error) {
	rows, err := r.queries.ListNotesByUser(ctx, dbgen.ListNotesByUserParams{UserID: db.PGUUID(userID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapNotes(rows), nil
}

func (r *postgresRepository) ListBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Note, error) {
	rows, err := r.queries.ListNotesBySource(ctx, dbgen.ListNotesBySourceParams{SourceID: db.PGUUID(sourceID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapNotes(rows), nil
}

func (r *postgresRepository) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Note, error) {
	rows, err := r.queries.ListPublicNotesByUser(ctx, dbgen.ListPublicNotesByUserParams{UserID: db.PGUUID(userID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapNotes(rows), nil
}

func (r *postgresRepository) ListPublicBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Note, error) {
	rows, err := r.queries.ListPublicNotesBySource(ctx, dbgen.ListPublicNotesBySourceParams{SourceID: db.PGUUID(sourceID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapNotes(rows), nil
}

func (r *postgresRepository) ListPublic(ctx context.Context, limit, offset int) ([]*Note, error) {
	rows, err := r.queries.ListPublicNotes(ctx, dbgen.ListPublicNotesParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapNotes(rows), nil
}

func (r *postgresRepository) Update(ctx context.Context, n *Note) (*Note, error) {
	annotations, err := json.Marshal(n.Annotations)
	if err != nil {
		return nil, err
	}
	tags, err := json.Marshal(n.Tags)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.UpdateNote(ctx, dbgen.UpdateNoteParams{
		ID:          db.PGUUID(n.ID),
		SourceID:    pgUUIDPtr(n.SourceID),
		Content:     n.Content,
		ContentType: string(n.ContentType),
		IsPublic:    db.PGBool(n.IsPublic),
		Annotations: annotations,
		Tags:        tags,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapNote(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteNote(ctx, db.PGUUID(id))
}

func (r *postgresRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountNotesByUser(ctx, db.PGUUID(userID))
}

func mapNotes(rows []dbgen.Note) []*Note {
	notes := make([]*Note, 0, len(rows))
	for _, row := range rows {
		notes = append(notes, mapNote(row))
	}
	return notes
}

func mapNote(row dbgen.Note) *Note {
	var annotations []string
	if len(row.Annotations) > 0 {
		_ = json.Unmarshal(row.Annotations, &annotations)
	}
	var tags []string
	if len(row.Tags) > 0 {
		_ = json.Unmarshal(row.Tags, &tags)
	}

	return &Note{
		ID:          db.UUID(row.ID),
		UserID:      db.UUID(row.UserID),
		SourceID:    uuidPtr(row.SourceID),
		Content:     row.Content,
		ContentType: ContentType(row.ContentType),
		IsPublic:    db.Bool(row.IsPublic),
		Annotations: annotations,
		Tags:        tags,
		CreatedAt:   db.Time(row.CreatedAt),
		UpdatedAt:   db.Time(row.UpdatedAt),
	}
}

func uuidPtr(value pgtype.UUID) *uuid.UUID {
	if !value.Valid {
		return nil
	}
	id := db.UUID(value)
	return &id
}

func pgUUIDPtr(value *uuid.UUID) pgtype.UUID {
	if value == nil {
		return pgtype.UUID{}
	}
	return db.PGUUID(*value)
}
