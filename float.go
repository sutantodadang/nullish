package nullish

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"strconv"

	"github.com/goccy/go-json"
)

type NullFloat struct {
	Float float64
	Valid bool
}

// Value method
func (nf NullFloat) Value() (driver.Value, error) {

	if !nf.Valid {
		return nil, nil
	}

	return nf.Float, nil
}

// Scan method
func (nf *NullFloat) Scan(value interface{}) error {

	if value == nil {
		nf.Float, nf.Valid = 0, false
		return nil
	}

	switch t := value.(type) {
	case float32:
		nf.Float, nf.Valid = float64(t), true

	case float64:
		nf.Float, nf.Valid = t, true

	case []byte:
		f, err := strconv.ParseFloat(string(t), 64)
		if err != nil {
			return errors.New("type assertion []byte to float64 is failed")
		}
		nf.Float, nf.Valid = f, true

	default:
		return errors.New("type assertion to float64 is failed")
	}

	return nil
}

// MarshalJSON method
func (nf NullFloat) MarshalJSON() ([]byte, error) {

	if !nf.Valid {
		return NullType, nil
	}

	return json.Marshal(nf.Float)
}

// UnmarshalJSON method
func (nf *NullFloat) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullType) {
		*nf = NullFloat{}
		return nil
	}

	var res float64

	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}

	*nf = NullFloat{Float: res, Valid: true}

	return nil
}
