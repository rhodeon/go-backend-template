package postgres

import "github.com/jackc/pgx/v5/pgtype"

// NewPgxText is a helper function to easily generate pgtype.Text from a string pointer.
// Non-empty strings return a valid text with the value.
// Nil pointers return an invalid text.
// Zero values return a result based on useZero. If at least one element of useZero is true a valid but empty text is returned.
// Otherwise, an invalid text is returned.
func NewPgxText(value *string, useZero ...bool) pgtype.Text {
	if value == nil {
		return pgtype.Text{
			Valid: false,
		}
	}

	if *value == "" {
		var validZero bool
		for _, uz := range useZero {
			validZero = validZero || uz
		}

		if !validZero {
			return pgtype.Text{
				Valid: false,
			}
		}
	}

	return pgtype.Text{
		String: *value,
		Valid:  true,
	}
}
