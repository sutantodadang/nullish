package nullish

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestNullInt_Value(t *testing.T) {
	ni := NewNullInt(42, true)
	got, err := ni.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Errorf("expected 42, got %v", got)
	}

	ni = NewNullInt(0, false)
	got, err = ni.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNullInt_Scan(t *testing.T) {
	var ni NullInt
	err := ni.Scan(int64(42))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ni.Int != 42 || !ni.Valid {
		t.Errorf("expected Int=42 Valid=true, got Int=%d Valid=%v", ni.Int, ni.Valid)
	}

	err = ni.Scan(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ni.Valid {
		t.Error("expected Valid=false for nil")
	}
}

func TestNullInt_JSON(t *testing.T) {
	ni := NewNullInt(42, true)
	data, err := json.Marshal(ni)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullInt
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if ni.Int != decoded.Int || ni.Valid != decoded.Valid {
		t.Errorf("roundtrip failed: expected %+v, got %+v", ni, decoded)
	}
}

func BenchmarkNullInt_Value(b *testing.B) {
	ni := NewNullInt(42, true)
	for i := 0; i < b.N; i++ {
		_, _ = ni.Value()
	}
}

func BenchmarkNullInt_Scan(b *testing.B) {
	val := int64(42)
	for i := 0; i < b.N; i++ {
		var ni NullInt
		_ = ni.Scan(val)
	}
}
