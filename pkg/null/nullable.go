package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Nullable[T any] struct {
	V     T
	Valid bool
}

func New[T any](value T, valid bool) Nullable[T] {
	return Nullable[T]{
		Valid: valid,
		V:     value,
	}
}

func (n Nullable[T]) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.V)
	}
	return ""
}

func (n Nullable[T]) IsZero() bool {
	return !n.Valid
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.V)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	var value *T
	n.Valid = false
	if err := json.Unmarshal(data, &value); err != nil {
		return nil
	}
	if value != nil {
		n.Valid = true
		n.V = *value
	}
	return nil
}

func (n *Nullable[T]) Scan(value any) error {
	if value == nil {
		n.V, n.Valid = *new(T), false
		return nil
	}
	n.Valid = true
	return nil
}

func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Value, nil
}

func ToSQL(n Nullable[bool]) sql.NullBool {
	return sql.NullBool{Valid: n.Valid, Bool: n.V}
}
