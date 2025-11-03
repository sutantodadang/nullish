package nullish

import (
	"testing"
	"time"

	"github.com/goccy/go-json"
)

func TestNullTime_Value(t *testing.T) {
	now := time.Now()
	nt := NewNullTime(now, true)
	got, err := nt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != now {
		t.Errorf("expected %v, got %v", now, got)
	}

	// Test invalid
	nt = NewNullTime(time.Time{}, false)
	got, err = nt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNullTime_Scan(t *testing.T) {
	now := time.Now()
	var nt NullTime

	// Scan time.Time
	err := nt.Scan(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !nt.Time.Equal(now) || !nt.Valid {
		t.Errorf("expected Time=%v Valid=true, got Time=%v Valid=%v", now, nt.Time, nt.Valid)
	}

	// Scan nil
	err = nt.Scan(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nt.Valid {
		t.Error("expected Valid=false for nil")
	}

	// Scan invalid type
	err = nt.Scan("not a time")
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestNullTime_MarshalJSON(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	nt := NewNullTime(now, true)

	data, err := nt.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `"` + now.Format(time.RFC3339Nano) + `"`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}

	// Test invalid
	nt = NewNullTime(time.Time{}, false)
	data, err = nt.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected null, got %s", string(data))
	}
}

func TestNullTime_UnmarshalJSON(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	jsonData := `"` + now.Format(time.RFC3339Nano) + `"`

	var nt NullTime
	err := nt.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !nt.Time.Equal(now) || !nt.Valid {
		t.Errorf("expected Time=%v Valid=true, got Time=%v Valid=%v", now, nt.Time, nt.Valid)
	}

	// Test null
	err = nt.UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if nt.Valid {
		t.Error("expected Valid=false for null")
	}
}

func TestNullTime_RoundTrip(t *testing.T) {
	now := time.Now()
	nt := NewNullTime(now, true)

	data, err := json.Marshal(nt)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullTime
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	// Use Equal for time comparison (handles location/monotonic)
	if !nt.Time.Equal(decoded.Time) || nt.Valid != decoded.Valid {
		t.Errorf("roundtrip failed: expected %+v, got %+v", nt, decoded)
	}
}

func BenchmarkNullTime_Value(b *testing.B) {
	nt := NewNullTime(time.Now(), true)
	for i := 0; i < b.N; i++ {
		_, _ = nt.Value()
	}
}

func BenchmarkNullTime_Scan(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		var nt NullTime
		_ = nt.Scan(now)
	}
}

func BenchmarkNullTime_MarshalJSON(b *testing.B) {
	nt := NewNullTime(time.Now(), true)
	for i := 0; i < b.N; i++ {
		_, _ = nt.MarshalJSON()
	}
}

func BenchmarkNullTime_UnmarshalJSON(b *testing.B) {
	now := time.Now()
	data := []byte(`"` + now.Format(time.RFC3339Nano) + `"`)
	for i := 0; i < b.N; i++ {
		var nt NullTime
		_ = nt.UnmarshalJSON(data)
	}
}
