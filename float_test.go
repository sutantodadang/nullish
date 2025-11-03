package nullish

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestNullFloat_Value(t *testing.T) {
	nf := NewNullFloat(3.14, true)
	got, err := nf.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3.14 {
		t.Errorf("expected 3.14, got %v", got)
	}
}

func TestNullFloat_Scan(t *testing.T) {
	var nf NullFloat
	err := nf.Scan(float64(3.14))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nf.Float != 3.14 || !nf.Valid {
		t.Errorf("expected Float=3.14 Valid=true, got Float=%f Valid=%v", nf.Float, nf.Valid)
	}
}

func TestNullFloat_JSON(t *testing.T) {
	nf := NewNullFloat(3.14, true)
	data, err := json.Marshal(nf)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullFloat
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if nf.Float != decoded.Float || nf.Valid != decoded.Valid {
		t.Errorf("roundtrip failed: expected %+v, got %+v", nf, decoded)
	}
}

func BenchmarkNullFloat_Value(b *testing.B) {
	nf := NewNullFloat(3.14, true)
	for i := 0; i < b.N; i++ {
		_, _ = nf.Value()
	}
}

func BenchmarkNullFloat_Scan(b *testing.B) {
	val := float64(3.14)
	for i := 0; i < b.N; i++ {
		var nf NullFloat
		_ = nf.Scan(val)
	}
}
