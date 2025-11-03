package nullish

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestNullArrObj_Value(t *testing.T) {
	arr := []map[string]interface{}{
		{"key": "value1"},
		{"key": "value2"},
	}
	na := NewNullArrObj(arr, true)

	got, err := na.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("expected non-nil value")
	}

	// Test invalid
	na = NewNullArrObj(nil, false)
	got, err = na.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNullArrObj_Scan(t *testing.T) {
	arr := []map[string]interface{}{{"key": "value"}}

	var na NullArrObj
	err := na.Scan(arr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !na.Valid {
		t.Error("expected Valid=true")
	}
	if len(na.ArrObj) != 1 {
		t.Errorf("expected length 1, got %d", len(na.ArrObj))
	}

	// Scan nil
	err = na.Scan(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if na.Valid {
		t.Error("expected Valid=false for nil")
	}

	// Scan invalid type
	err = na.Scan("not an array")
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestNullArrObj_MarshalJSON(t *testing.T) {
	arr := []map[string]interface{}{{"key": "value"}}
	na := NewNullArrObj(arr, true)

	data, err := na.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if string(data) != `[{"key":"value"}]` {
		t.Errorf("expected [{\"key\":\"value\"}], got %s", string(data))
	}

	// Test invalid
	na = NewNullArrObj(nil, false)
	data, err = na.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected null, got %s", string(data))
	}
}

func TestNullArrObj_UnmarshalJSON(t *testing.T) {
	jsonData := `[{"key":"value1"},{"key":"value2"}]`

	var na NullArrObj
	err := na.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !na.Valid {
		t.Error("expected Valid=true")
	}
	if len(na.ArrObj) != 2 {
		t.Errorf("expected length 2, got %d", len(na.ArrObj))
	}
	if na.ArrObj[0]["key"] != "value1" {
		t.Errorf("expected first element key=value1, got %v", na.ArrObj[0]["key"])
	}

	// Test null
	err = na.UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if na.Valid {
		t.Error("expected Valid=false for null")
	}
}

func TestNullArrObj_RoundTrip(t *testing.T) {
	arr := []map[string]interface{}{
		{"name": "test1", "value": float64(1)},
		{"name": "test2", "value": float64(2)},
	}
	na := NewNullArrObj(arr, true)

	data, err := json.Marshal(na)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullArrObj
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if na.Valid != decoded.Valid {
		t.Errorf("Valid mismatch: expected %v, got %v", na.Valid, decoded.Valid)
	}
	if len(na.ArrObj) != len(decoded.ArrObj) {
		t.Errorf("length mismatch: expected %d, got %d", len(na.ArrObj), len(decoded.ArrObj))
	}
}

func BenchmarkNullArrObj_Value(b *testing.B) {
	na := NewNullArrObj([]map[string]interface{}{{"key": "value"}}, true)
	for i := 0; i < b.N; i++ {
		_, _ = na.Value()
	}
}

func BenchmarkNullArrObj_Scan(b *testing.B) {
	arr := []map[string]interface{}{{"key": "value"}}
	for i := 0; i < b.N; i++ {
		var na NullArrObj
		_ = na.Scan(arr)
	}
}

func BenchmarkNullArrObj_MarshalJSON(b *testing.B) {
	na := NewNullArrObj([]map[string]interface{}{{"key": "value"}}, true)
	for i := 0; i < b.N; i++ {
		_, _ = na.MarshalJSON()
	}
}

func BenchmarkNullArrObj_UnmarshalJSON(b *testing.B) {
	data := []byte(`[{"key":"value"}]`)
	for i := 0; i < b.N; i++ {
		var na NullArrObj
		_ = na.UnmarshalJSON(data)
	}
}
