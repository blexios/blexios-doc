package parser

import (
	"testing"
)

func TestParseValidArrayOfObjects(t *testing.T) {
	data := []byte(`[
        {"id": "1", "name": "John"},
        {"id": "2", "name": "Jane"}
    ]`)

	rows, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON returned error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0]["name"] != "John" {
		t.Fatalf("expected first row name to be John, got %v", rows[0]["name"])
	}
}

func TestParseEmptyArray(t *testing.T) {
	_, err := ParseJSON([]byte(`[]`))
	if err == nil {
		t.Fatal("expected empty array parse to fail")
	}
}

func TestParseInvalidJSON(t *testing.T) {
	_, err := ParseJSON([]byte(`[{"id": 1`))
	if err == nil {
		t.Fatal("expected invalid JSON to fail")
	}
}

func TestDetectColumns(t *testing.T) {
	rows := []map[string]any{
		{"id": "1", "name": "John", "active": true},
		{"id": "2", "name": "Jane", "active": false},
	}

	cols, err := DetectColumns(rows)
	if err != nil {
		t.Fatalf("DetectColumns returned error: %v", err)
	}
	if len(cols) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(cols))
	}
	if cols[0].Key != "id" {
		t.Fatalf("expected first column key to be id, got %s", cols[0].Key)
	}
}
