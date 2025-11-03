package nullish

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestNullBool_Value(t *testing.T) {
	nb := NewNullBool(true, true)
	got, err := nb.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}

	nb = NewNullBool(false, false)
	got, err = nb.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNullBool_Scan(t *testing.T) {
	var nb NullBool
	err := nb.Scan(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !nb.Bool || !nb.Valid {
		t.Errorf("expected Bool=true Valid=true, got Bool=%v Valid=%v", nb.Bool, nb.Valid)
	}

	err = nb.Scan(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nb.Valid {
		t.Error("expected Valid=false for nil")
	}
}

func TestNullBool_JSON(t *testing.T) {
	nb := NewNullBool(true, true)
	data, err := json.Marshal(nb)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullBool
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if nb.Bool != decoded.Bool || nb.Valid != decoded.Valid {
		t.Errorf("roundtrip failed: expected %+v, got %+v", nb, decoded)
	}
}

func BenchmarkNullBool_Value(b *testing.B) {
	nb := NewNullBool(true, true)
	for i := 0; i < b.N; i++ {
		_, _ = nb.Value()
	}
}

func BenchmarkNullBool_Scan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var nb NullBool
		_ = nb.Scan(true)
	}
}
