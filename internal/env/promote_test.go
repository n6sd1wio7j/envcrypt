package env

import (
	"testing"
)

func TestPromote_AllKeys_NoConflict(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "1"}, {Key: "BAR", Value: "2"}}
	dst := []Entry{{Key: "BAZ", Value: "3"}}

	out, res, err := Promote(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries in output, got %d", len(out))
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "new"}}
	dst := []Entry{{Key: "FOO", Value: "old"}}

	out, res, err := Promote(src, dst, PromoteOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected FOO to be skipped, got %v", res.Skipped)
	}
	if out[0].Value != "old" {
		t.Errorf("expected value to remain 'old', got %q", out[0].Value)
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "new"}}
	dst := []Entry{{Key: "FOO", Value: "old"}}

	out, res, err := Promote(src, dst, PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(res.Overwritten))
	}
	if out[0].Value != "new" {
		t.Errorf("expected value 'new', got %q", out[0].Value)
	}
}

func TestPromote_KeyFilter(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "1"}, {Key: "BAR", Value: "2"}}
	dst := []Entry{}

	out, res, err := Promote(src, dst, PromoteOptions{Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "FOO" {
		t.Errorf("expected only FOO promoted, got %v", res.Promoted)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry in output, got %d", len(out))
	}
}

func TestPromote_DryRun_DoesNotModify(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "1"}}
	dst := []Entry{{Key: "BAR", Value: "2"}}

	out, res, err := Promote(src, dst, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("dry-run should not add entries, got %d", len(out))
	}
	if len(res.Promoted) != 1 {
		t.Errorf("expected 1 in promoted report, got %d", len(res.Promoted))
	}
}

func TestPromote_NilSource_ReturnsError(t *testing.T) {
	_, _, err := Promote(nil, []Entry{}, PromoteOptions{})
	if err == nil {
		t.Error("expected error for nil source, got nil")
	}
}
