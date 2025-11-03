package nullish

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestNullString_Value(t *testing.T) {
	ns := NewNullString("test", true)
	got, err := ns.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "test" {
		t.Errorf("expected 'test', got %v", got)
	}

	ns = NewNullString("", false)
	got, err = ns.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNullString_Scan(t *testing.T) {
	var ns NullString
	err := ns.Scan("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ns.String != "test" || !ns.Valid {
		t.Errorf("expected String='test' Valid=true, got String=%q Valid=%v", ns.String, ns.Valid)
	}

	err = ns.Scan(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ns.Valid {
		t.Error("expected Valid=false for nil")
	}
}

func TestNullString_JSON(t *testing.T) {
	ns := NewNullString("test", true)
	data, err := json.Marshal(ns)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullString
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if ns.String != decoded.String || ns.Valid != decoded.Valid {
		t.Errorf("roundtrip failed: expected %+v, got %+v", ns, decoded)
	}
}

func BenchmarkNullString_Value(b *testing.B) {
	ns := NewNullString("benchmark", true)
	for i := 0; i < b.N; i++ {
		_, _ = ns.Value()
	}
}

func BenchmarkNullString_Scan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ns NullString
		_ = ns.Scan("benchmark")
	}
}
