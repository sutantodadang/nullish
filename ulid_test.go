package nullish

import (
	"math/rand"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/oklog/ulid/v2"
)

func TestNullULID_Value(t *testing.T) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	nl := NewNullULID(id, true)

	got, err := nl.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("expected non-nil value")
	}

	// Test invalid
	nl = NewNullULID(ulid.ULID{}, false)
	got, err = nl.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNullULID_Scan(t *testing.T) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	var nl NullULID

	// Scan ULID string
	err := nl.Scan(id.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nl.ULID.String() != id.String() || !nl.Valid {
		t.Errorf("expected ULID=%v Valid=true, got ULID=%v Valid=%v", id, nl.ULID, nl.Valid)
	}

	// Scan ULID bytes
	idBytes := id[:]
	err = nl.Scan(idBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nl.ULID != id || !nl.Valid {
		t.Errorf("expected ULID=%v Valid=true, got ULID=%v Valid=%v", id, nl.ULID, nl.Valid)
	}

	// Scan nil
	err = nl.Scan(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nl.Valid {
		t.Error("expected Valid=false for nil")
	}

	// Scan invalid type
	err = nl.Scan(123)
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestNullULID_MarshalJSON(t *testing.T) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	nl := NewNullULID(id, true)

	data, err := nl.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `"` + id.String() + `"`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}

	// Test invalid
	nl = NewNullULID(ulid.ULID{}, false)
	data, err = nl.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected null, got %s", string(data))
	}
}

func TestNullULID_UnmarshalJSON(t *testing.T) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	jsonData := `"` + id.String() + `"`

	var nl NullULID
	err := nl.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if nl.ULID.String() != id.String() || !nl.Valid {
		t.Errorf("expected ULID=%v Valid=true, got ULID=%v Valid=%v", id, nl.ULID, nl.Valid)
	}

	// Test null
	err = nl.UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if nl.Valid {
		t.Error("expected Valid=false for null")
	}
}

func TestNullULID_RoundTrip(t *testing.T) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	nl := NewNullULID(id, true)

	data, err := json.Marshal(nl)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded NullULID
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if nl.ULID != decoded.ULID || nl.Valid != decoded.Valid {
		t.Errorf("roundtrip failed: expected %+v, got %+v", nl, decoded)
	}
}

func BenchmarkNullULID_Value(b *testing.B) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	nl := NewNullULID(ulid.MustNew(ulid.Timestamp(time.Now()), entropy), true)
	for i := 0; i < b.N; i++ {
		_, _ = nl.Value()
	}
}

func BenchmarkNullULID_Scan_String(b *testing.B) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	for i := 0; i < b.N; i++ {
		var nl NullULID
		_ = nl.Scan(id)
	}
}

func BenchmarkNullULID_Scan_Bytes(b *testing.B) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	testULID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	id := testULID[:]
	for i := 0; i < b.N; i++ {
		var nl NullULID
		_ = nl.Scan(id)
	}
}

func BenchmarkNullULID_MarshalJSON(b *testing.B) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	nl := NewNullULID(ulid.MustNew(ulid.Timestamp(time.Now()), entropy), true)
	for i := 0; i < b.N; i++ {
		_, _ = nl.MarshalJSON()
	}
}

func BenchmarkNullULID_UnmarshalJSON(b *testing.B) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	data := []byte(`"` + id.String() + `"`)
	for i := 0; i < b.N; i++ {
		var nl NullULID
		_ = nl.UnmarshalJSON(data)
	}
}
