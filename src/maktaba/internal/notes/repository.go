package notes

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zizouhuweidi/maktaba/internal/db"
)

type postgresRepository struct {
	db *db.DB
}

// NewPostgresRepository creates a new postgres repository for notes
func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{
		db: d,
	}
}

func (r *postgresRepository) Create(ctx context.Context, n *Note) (*Note, error) {
	row, err := r.db.Queries.CreateNote(ctx, db.CreateNoteParams{
		UserID:      n.UserID,
		SourceID:    r.toUUID(n.SourceID),
		Content:     n.Content,
		ContentType: string(n.ContentType),
		IsPublic:    r.toBool(n.IsPublic),
		Annotations: r.toRawJSON(n.Annotations),
		Tags:        r.toRawJSON(n.Tags),
	})
	if err != nil {
		return nil, err
	}

	return r.mapToNote(row), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Note, error) {
	row, err := r.db.Queries.GetNote(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapToNote(row), nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Note, error) {
	rows, err := r.db.Queries.ListNotesByUser(ctx, db.ListNotesByUserParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	notes := make([]*Note, len(rows))
	for i, row := range rows {
		notes[i] = r.mapToNote(row)
	}
	return notes, nil
}

func (r *postgresRepository) ListBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Note, error) {
	rows, err := r.db.Queries.ListNotesBySource(ctx, db.ListNotesBySourceParams{
		SourceID: r.toUUID(&sourceID),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	notes := make([]*Note, len(rows))
	for i, row := range rows {
		notes[i] = r.mapToNote(row)
	}
	return notes, nil
}

func (r *postgresRepository) ListPublic(ctx context.Context, limit, offset int) ([]*Note, error) {
	rows, err := r.db.Queries.ListPublicNotes(ctx, db.ListPublicNotesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	notes := make([]*Note, len(rows))
	for i, row := range rows {
		notes[i] = r.mapToNote(row)
	}
	return notes, nil
}

func (r *postgresRepository) Update(ctx context.Context, n *Note) (*Note, error) {
	row, err := r.db.Queries.UpdateNote(ctx, db.UpdateNoteParams{
		ID:          n.ID,
		Content:     r.toText(r.PtrString(n.Content)),
		ContentType: r.toText(r.PtrString(string(n.ContentType))),
		IsPublic:    r.toBool(n.IsPublic),
		Annotations: r.toRawJSON(n.Annotations),
		Tags:        r.toRawJSON(n.Tags),
	})
	if err != nil {
		return nil, err
	}

	return r.mapToNote(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Queries.DeleteNote(ctx, id)
}

func (r *postgresRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.db.Queries.CountNotesByUser(ctx, userID)
}

func (r *postgresRepository) mapToNote(row db.Note) *Note {
	var annotations []string
	if row.Annotations != nil {
		json.Unmarshal(row.Annotations, &annotations)
	}

	var tags []string
	if row.Tags != nil {
		json.Unmarshal(row.Tags, &tags)
	}

	return &Note{
		ID:          row.ID,
		UserID:      row.UserID,
		SourceID:    r.fromUUID(row.SourceID),
		Content:     row.Content,
		ContentType: ContentType(row.ContentType),
		IsPublic:    row.IsPublic.Bool,
		Annotations: annotations,
		Tags:        tags,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

// Helper methods for pgtype conversions
func (r *postgresRepository) toText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func (r *postgresRepository) fromText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func (r *postgresRepository) toUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func (r *postgresRepository) fromUUID(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	res := uuid.UUID(u.Bytes)
	return &res
}

func (r *postgresRepository) toBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func (r *postgresRepository) toRawJSON(data interface{}) []byte {
	if data == nil {
		return nil
	}
	b, _ := json.Marshal(data)
	return b
}

func (r *postgresRepository) PtrString(s string) *string {
	return &s
}
