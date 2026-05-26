package reviews

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapCreateError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "duplicate", err: &pgconn.PgError{Code: "23505"}, want: ErrReviewExists},
		{name: "missing source", err: &pgconn.PgError{Code: "23503", ConstraintName: "reviews_source_id_fkey"}, want: ErrSourceNotFound},
		{name: "other foreign key", err: &pgconn.PgError{Code: "23503", ConstraintName: "reviews_user_id_fkey"}, want: ErrReviewConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mapCreateError(tt.err); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}
