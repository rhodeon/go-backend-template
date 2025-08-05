package postgres

import (
	"testing"

	"github.com/rhodeon/go-backend-template/utils/typeutils"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestNewPgxText(t *testing.T) {
	type args struct {
		value   *string
		useZero []bool
	}

	testCases := map[string]struct {
		args args
		want pgtype.Text
	}{
		"Valid value": {
			args: args{
				value: typeutils.Ptr("John Doe"),
			},
			want: pgtype.Text{
				String: "John Doe",
				Valid:  true,
			},
		},

		"Invalid zero value": {
			args: args{
				value: typeutils.Ptr(""),
			},
			want: pgtype.Text{
				String: "",
				Valid:  false,
			},
		},

		"Valid zero value": {
			args: args{
				value:   typeutils.Ptr(""),
				useZero: []bool{true},
			},
			want: pgtype.Text{
				String: "",
				Valid:  true,
			},
		},

		"Nil value": {
			args: args{
				value:   nil,
				useZero: []bool{true},
			},
			want: pgtype.Text{
				String: "",
				Valid:  false,
			},
		},

		"Zero value with multiple useZero one of which is true": {
			args: args{
				value:   typeutils.Ptr(""),
				useZero: []bool{false, true, false},
			},
			want: pgtype.Text{
				String: "",
				Valid:  true,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := NewPgxText(tc.args.value, tc.args.useZero...)
			assert.Equal(t, tc.want, got)
		})
	}
}
