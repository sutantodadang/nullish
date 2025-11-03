package nullish

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/goccy/go-json"
)

type NullArrObj struct {
	ArrObj []map[string]interface{}
	Valid  bool
}

// Value method
func (na NullArrObj) Value() (driver.Value, error) {

	if !na.Valid {
		return nil, nil
	}

	return json.Marshal(na.ArrObj)
}

// Scan method
func (na *NullArrObj) Scan(value interface{}) error {

	if value == nil {
		na.ArrObj, na.Valid = []map[string]interface{}{}, false
		return nil
	}

	b, ok := value.([]map[string]interface{})
	if !ok {
		return errors.New("type assertion to array object is failed")
	}

	na.ArrObj, na.Valid = b, true

	return nil
}

// MarshalJSON method
func (na NullArrObj) MarshalJSON() ([]byte, error) {

	if !na.Valid {
		return NullType, nil
	}

	return json.Marshal(na.ArrObj)
}

// UnmarshalJSON method
func (na *NullArrObj) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullType) {
		*na = NullArrObj{}
		return nil
	}

	var res []map[string]interface{}

	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}

	*na = NullArrObj{ArrObj: res, Valid: true}

	return nil
}
