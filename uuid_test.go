package nullish

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

func TestNullUUID_Value(t *testing.T) {
	id := uuid.New()
	nu := NewNullUUID(id, true)

	got, err := nu.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("expected non-nil value")
	}

	// Test invalid
	nu = NewNullUUID(uuid.UUID{}, false)
	got, err = nu.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNullUUID_Scan(t *testing.T) {
	id := uuid.New()
	var nu NullUUID

	// Scan UUID string
	err := nu.Scan(id.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nu.UUID != id || !nu.Valid {
		t.Errorf("expected UUID=%v Valid=true, got UUID=%v Valid=%v", id, nu.UUID, nu.Valid)
	}

	// Scan UUID bytes
	idBytes := id[:]
	err = nu.Scan(idBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nu.UUID != id || !nu.Valid {
		t.Errorf("expected UUID=%v Valid=true, got UUID=%v Valid=%v", id, nu.UUID, nu.Valid)
	}

	// Scan nil
	err = nu.Scan(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nu.Valid {
		t.Error("expected Valid=false for nil")
	}

	// Scan invalid type
	err = nu.Scan(123)
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestNullUUID_MarshalJSON(t *testing.T) {
	id := uuid.New()
	nu := NewNullUUID(id, true)

	data, err := nu.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `"` + id.String() + `"`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}

	// Test invalid
	nu = NewNullUUID(uuid.UUID{}, false)
	data, err = nu.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected null, got %s", string(data))
	}
}

func TestNullUUID_UnmarshalJSON(t *testing.T) {
	id := uuid.New()
	jsonData := `"` + id.String() + `"`

	var nu NullUUID
	err := nu.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if nu.UUID != id || !nu.Valid {
		t.Errorf("expected UUID=%v Valid=true, got UUID=%v Valid=%v", id, nu.UUID, nu.Valid)
	}

	// Test null
	err = nu.UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if nu.Valid {
		t.Error("expected Valid=false for null")
	}
}

func TestNullUUID_RoundTrip(t *testing.T) {
	id := uuid.New()
	nu := NewNullUUID(id, true)

	data, err := json.Marshal(nu)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullUUID
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if nu.UUID != decoded.UUID || nu.Valid != decoded.Valid {
		t.Errorf("roundtrip failed: expected %+v, got %+v", nu, decoded)
	}
}

func BenchmarkNullUUID_Value(b *testing.B) {
	nu := NewNullUUID(uuid.New(), true)
	for i := 0; i < b.N; i++ {
		_, _ = nu.Value()
	}
}

func BenchmarkNullUUID_Scan_String(b *testing.B) {
	id := uuid.New().String()
	for i := 0; i < b.N; i++ {
		var nu NullUUID
		_ = nu.Scan(id)
	}
}

func BenchmarkNullUUID_Scan_Bytes(b *testing.B) {
	testUUID := uuid.New()
	id := testUUID[:]
	for i := 0; i < b.N; i++ {
		var nu NullUUID
		_ = nu.Scan(id)
	}
}

func BenchmarkNullUUID_MarshalJSON(b *testing.B) {
	nu := NewNullUUID(uuid.New(), true)
	for i := 0; i < b.N; i++ {
		_, _ = nu.MarshalJSON()
	}
}

func BenchmarkNullUUID_UnmarshalJSON(b *testing.B) {
	id := uuid.New()
	data := []byte(`"` + id.String() + `"`)
	for i := 0; i < b.N; i++ {
		var nu NullUUID
		_ = nu.UnmarshalJSON(data)
	}
}
