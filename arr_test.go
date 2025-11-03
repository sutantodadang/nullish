package nullish

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestNullArr_Value(t *testing.T) {
	arr := []interface{}{"a", "b", float64(1), float64(2)}
	na := NewNullArr(arr, true)

	got, err := na.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("expected non-nil value")
	}

	// Test invalid
	na = NewNullArr(nil, false)
	got, err = na.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNullArr_Scan(t *testing.T) {
	arr := []interface{}{"a", "b"}

	var na NullArr
	err := na.Scan(arr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !na.Valid {
		t.Error("expected Valid=true")
	}
	if len(na.Arr) != 2 {
		t.Errorf("expected length 2, got %d", len(na.Arr))
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

func TestNullArr_MarshalJSON(t *testing.T) {
	arr := []interface{}{"a", "b"}
	na := NewNullArr(arr, true)

	data, err := na.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if string(data) != `["a","b"]` {
		t.Errorf("expected [\"a\",\"b\"], got %s", string(data))
	}

	// Test invalid
	na = NewNullArr(nil, false)
	data, err = na.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected null, got %s", string(data))
	}
}

func TestNullArr_UnmarshalJSON(t *testing.T) {
	jsonData := `["a","b",1,2]`

	var na NullArr
	err := na.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !na.Valid {
		t.Error("expected Valid=true")
	}
	if len(na.Arr) != 4 {
		t.Errorf("expected length 4, got %d", len(na.Arr))
	}
	if na.Arr[0] != "a" {
		t.Errorf("expected first element 'a', got %v", na.Arr[0])
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

func TestNullArr_RoundTrip(t *testing.T) {
	arr := []interface{}{"test", float64(123)}
	na := NewNullArr(arr, true)

	data, err := json.Marshal(na)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullArr
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if na.Valid != decoded.Valid {
		t.Errorf("Valid mismatch: expected %v, got %v", na.Valid, decoded.Valid)
	}
	if len(na.Arr) != len(decoded.Arr) {
		t.Errorf("length mismatch: expected %d, got %d", len(na.Arr), len(decoded.Arr))
	}
}

func BenchmarkNullArr_Value(b *testing.B) {
	na := NewNullArr([]interface{}{"a", "b"}, true)
	for i := 0; i < b.N; i++ {
		_, _ = na.Value()
	}
}

func BenchmarkNullArr_Scan(b *testing.B) {
	arr := []interface{}{"a", "b"}
	for i := 0; i < b.N; i++ {
		var na NullArr
		_ = na.Scan(arr)
	}
}

func BenchmarkNullArr_MarshalJSON(b *testing.B) {
	na := NewNullArr([]interface{}{"a", "b"}, true)
	for i := 0; i < b.N; i++ {
		_, _ = na.MarshalJSON()
	}
}

func BenchmarkNullArr_UnmarshalJSON(b *testing.B) {
	data := []byte(`["a","b"]`)
	for i := 0; i < b.N; i++ {
		var na NullArr
		_ = na.UnmarshalJSON(data)
	}
}
