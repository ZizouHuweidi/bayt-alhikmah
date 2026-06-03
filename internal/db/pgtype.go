package db

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func PGUUID(id uuid.UUID) pgtype.UUID {
	if id == uuid.Nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: [16]byte(id), Valid: true}
}

func UUID(value pgtype.UUID) uuid.UUID {
	return uuid.UUID(value.Bytes)
}

func PGText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func PGTextString(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}

func StringPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func PGBool(value bool) pgtype.Bool {
	return pgtype.Bool{Bool: value, Valid: true}
}

func Bool(value pgtype.Bool) bool {
	return value.Valid && value.Bool
}

func PGInt4Ptr(value *int) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: int32(*value), Valid: true}
}

func IntPtr(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	v := int(value.Int32)
	return &v
}

func PGTimestamptzPtr(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func PGTimestamptz(value time.Time) pgtype.Timestamptz {
	if value.IsZero() {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: value, Valid: true}
}

func Time(value pgtype.Timestamptz) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}

func TimePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}
