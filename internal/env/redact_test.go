package env

import (
	"testing"
)

func TestRedact_DefaultPatterns(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "AUTH_TOKEN", Value: "tok-xyz"},
		{Key: "PORT", Value: "8080"},
	}

	result := Redact(entries, RedactOptions{})

	if result[0].Value != "myapp" {
		t.Errorf("APP_NAME should not be redacted, got %q", result[0].Value)
	}
	if result[1].Value != "***" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result[1].Value)
	}
	if result[2].Value != "***" {
		t.Errorf("API_KEY should be redacted, got %q", result[2].Value)
	}
	if result[3].Value != "***" {
		t.Errorf("AUTH_TOKEN should be redacted, got %q", result[3].Value)
	}
	if result[4].Value != "8080" {
		t.Errorf("PORT should not be redacted, got %q", result[4].Value)
	}
}

func TestRedact_ExplicitKeys(t *testing.T) {
	entries := []Entry{
		{Key: "CUSTOM_FIELD", Value: "sensitive"},
		{Key: "NORMAL", Value: "public"},
	}

	result := Redact(entries, RedactOptions{Keys: []string{"CUSTOM_FIELD"}})

	if result[0].Value != "***" {
		t.Errorf("CUSTOM_FIELD should be redacted, got %q", result[0].Value)
	}
	if result[1].Value != "public" {
		t.Errorf("NORMAL should not be redacted, got %q", result[1].Value)
	}
}

func TestRedact_CustomMask(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "secret"},
	}

	result := Redact(entries, RedactOptions{Mask: "[REDACTED]"})

	if result[0].Value != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", result[0].Value)
	}
}

func TestRedact_CustomPattern(t *testing.T) {
	entries := []Entry{
		{Key: "STRIPE_LIVE_KEY", Value: "sk_live_abc"},
		{Key: "SAFE", Value: "value"},
	}

	result := Redact(entries, RedactOptions{Patterns: []string{`(?i)stripe`}})

	if result[0].Value != "***" {
		t.Errorf("STRIPE_LIVE_KEY should be redacted, got %q", result[0].Value)
	}
	if result[1].Value != "value" {
		t.Errorf("SAFE should not be redacted, got %q", result[1].Value)
	}
}

func TestRedactMap_Basic(t *testing.T) {
	m := map[string]string{
		"DB_SECRET": "topsecret",
		"HOST":      "localhost",
	}

	result := RedactMap(m, RedactOptions{})

	if result["DB_SECRET"] != "***" {
		t.Errorf("DB_SECRET should be redacted, got %q", result["DB_SECRET"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("HOST should not be redacted, got %q", result["HOST"])
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "original"},
	}

	_ = Redact(entries, RedactOptions{})

	if entries[0].Value != "original" {
		t.Errorf("original entries should not be mutated, got %q", entries[0].Value)
	}
}
