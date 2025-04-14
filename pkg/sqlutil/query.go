package sqlutil

import (
	"database/sql"
)

type ScanFunc func(dest ...any) error

func ScanRows[K any](rows *sql.Rows, scan func(sf ScanFunc, k K) error) ([]K, error) {
	defer rows.Close()
	var items []K
	for rows.Next() {
		var k K
		err := scan(rows.Scan, k)
		if err != nil {
			return nil, err
		}
		items = append(items, k)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
