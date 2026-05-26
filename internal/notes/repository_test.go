package notes

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapCreateErrorMissingSource(t *testing.T) {
	err := mapCreateError(&pgconn.PgError{Code: "23503", ConstraintName: "notes_source_id_fkey"})
	if !errors.Is(err, ErrSourceNotFound) {
		t.Fatalf("error = %v, want %v", err, ErrSourceNotFound)
	}
}
