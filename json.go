package nullish

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/goccy/go-json"
)

type NullJSON struct {
	Json  json.RawMessage
	Valid bool
}

// Value method
func (nj NullJSON) Value() (driver.Value, error) {

	if !nj.Valid {
		return nil, nil
	}

	return json.Marshal(nj.Json)
}

// Scan method
func (nj *NullJSON) Scan(value interface{}) error {

	if value == nil {
		nj.Json, nj.Valid = json.RawMessage{}, false
		return nil
	}

	switch t := value.(type) {
	case string:
		nj.Json, nj.Valid = []byte(t), true

	case []byte:
		if len(t) == 0 {
			nj.Json, nj.Valid = NullType, true
		} else {
			nj.Json, nj.Valid = t, true
		}

	default:
		return errors.New("invalid type json")
	}

	return nil
}

// MarshalJSON method
func (nj NullJSON) MarshalJSON() ([]byte, error) {

	if !nj.Valid {
		return NullType, nil
	}

	return json.Marshal(nj.Json)
}

// UnmarshalJSON method
func (nj *NullJSON) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullType) {
		*nj = NullJSON{}
		return nil
	}

	*nj = NullJSON{Json: data, Valid: true}

	return nil
}
