package formatter

import "testing"

func TestFormatNullValue(t *testing.T) {
	if got := FormatValue(nil, "string"); got != "—" {
		t.Fatalf("expected null to be formatted as em dash, got %q", got)
	}
}

func TestFormatDateValue(t *testing.T) {
	if got := FormatValue("2026-09-21", "date"); got != "21/09/2026" {
		t.Fatalf("expected date formatted to 21/09/2026, got %q", got)
	}
}

func TestFormatTimeValue(t *testing.T) {
	if got := FormatValue("19:37:00", "time"); got != "19:37" {
		t.Fatalf("expected time formatted to 19:37, got %q", got)
	}
}
