package nullish

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestNullObj_Value(t *testing.T) {
	obj := map[string]interface{}{"key": "value", "num": float64(42)}
	no := NewNullObj(obj, true)

	got, err := no.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("expected non-nil value")
	}

	// Test invalid
	no = NewNullObj(nil, false)
	got, err = no.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNullObj_Scan(t *testing.T) {
	obj := map[string]interface{}{"key": "value"}

	var no NullObj
	err := no.Scan(obj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !no.Valid {
		t.Error("expected Valid=true")
	}
	if no.Obj["key"] != "value" {
		t.Errorf("expected key=value, got %v", no.Obj["key"])
	}

	// Scan nil
	err = no.Scan(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if no.Valid {
		t.Error("expected Valid=false for nil")
	}

	// Scan invalid type
	err = no.Scan("not an object")
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestNullObj_MarshalJSON(t *testing.T) {
	obj := map[string]interface{}{"key": "value"}
	no := NewNullObj(obj, true)

	data, err := no.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if string(data) != `{"key":"value"}` {
		t.Errorf("expected {\"key\":\"value\"}, got %s", string(data))
	}

	// Test invalid
	no = NewNullObj(nil, false)
	data, err = no.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected null, got %s", string(data))
	}
}

func TestNullObj_UnmarshalJSON(t *testing.T) {
	jsonData := `{"key":"value","num":42}`

	var no NullObj
	err := no.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !no.Valid {
		t.Error("expected Valid=true")
	}
	if no.Obj["key"] != "value" {
		t.Errorf("expected key=value, got %v", no.Obj["key"])
	}
	if no.Obj["num"] != float64(42) {
		t.Errorf("expected num=42, got %v", no.Obj["num"])
	}

	// Test null
	err = no.UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if no.Valid {
		t.Error("expected Valid=false for null")
	}
}

func TestNullObj_RoundTrip(t *testing.T) {
	obj := map[string]interface{}{"name": "test", "value": float64(123)}
	no := NewNullObj(obj, true)

	data, err := json.Marshal(no)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullObj
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if no.Valid != decoded.Valid {
		t.Errorf("Valid mismatch: expected %v, got %v", no.Valid, decoded.Valid)
	}
	if no.Obj["name"] != decoded.Obj["name"] {
		t.Errorf("name mismatch: expected %v, got %v", no.Obj["name"], decoded.Obj["name"])
	}
}

func BenchmarkNullObj_Value(b *testing.B) {
	no := NewNullObj(map[string]interface{}{"key": "value"}, true)
	for i := 0; i < b.N; i++ {
		_, _ = no.Value()
	}
}

func BenchmarkNullObj_Scan(b *testing.B) {
	obj := map[string]interface{}{"key": "value"}
	for i := 0; i < b.N; i++ {
		var no NullObj
		_ = no.Scan(obj)
	}
}

func BenchmarkNullObj_MarshalJSON(b *testing.B) {
	no := NewNullObj(map[string]interface{}{"key": "value"}, true)
	for i := 0; i < b.N; i++ {
		_, _ = no.MarshalJSON()
	}
}

func BenchmarkNullObj_UnmarshalJSON(b *testing.B) {
	data := []byte(`{"key":"value"}`)
	for i := 0; i < b.N; i++ {
		var no NullObj
		_ = no.UnmarshalJSON(data)
	}
}
