package pgxutil

import (
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ToArray[T comparable](rows pgx.Rows, err error) ([]T, error) {
	if err != nil {
		return nil, fmt.Errorf("query:%w", err)
	}

	var out []T

	for rows.Next() {
		var row T

		err := rows.Scan(&row)
		if err != nil {
			return nil, fmt.Errorf("scan:%w", err)
		}

		out = append(out, row)
	}

	return out, nil
}
