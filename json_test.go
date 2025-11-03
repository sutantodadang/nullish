package nullish

import (
	"testing"

	"github.com/goccy/go-json"
)

func Test_NullJSON(t *testing.T) {

	testData := []map[string]interface{}{
		{
			"foo": "bar",
		},
	}

	b, err := json.Marshal(&testData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jsonNull := NewNullJSON(b, true)

	if !jsonNull.Valid {
		t.Error("expected Valid to be true")
	}

	var newTestData []map[string]interface{}

	nb, err := jsonNull.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = json.Unmarshal(nb, &newTestData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(testData) != len(newTestData) {
		t.Errorf("expected same length, got %d vs %d", len(testData), len(newTestData))
	}
	if testData[0]["foo"] != newTestData[0]["foo"] {
		t.Errorf("expected %v, got %v", testData, newTestData)
	}
}
