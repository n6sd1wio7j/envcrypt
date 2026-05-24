package env

import (
	"strings"
	"testing"
)

func TestCompare_Identical(t *testing.T) {
	a := []Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}
	b := []Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}
	r := Compare(a, b)
	if len(r.Identical) != 2 {
		t.Errorf("expected 2 identical, got %d", len(r.Identical))
	}
	if len(r.Changed) != 0 || len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 {
		t.Error("expected no differences")
	}
}

func TestCompare_Changed(t *testing.T) {
	a := []Entry{{Key: "FOO", Value: "old"}}
	b := []Entry{{Key: "FOO", Value: "new"}}
	r := Compare(a, b)
	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(r.Changed))
	}
	if r.Changed[0].ValueA != "old" || r.Changed[0].ValueB != "new" {
		t.Errorf("unexpected changed values: %+v", r.Changed[0])
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	a := []Entry{{Key: "ONLY_A", Value: "1"}, {Key: "SHARED", Value: "x"}}
	b := []Entry{{Key: "SHARED", Value: "x"}}
	r := Compare(a, b)
	if len(r.OnlyInA) != 1 || r.OnlyInA[0].Key != "ONLY_A" {
		t.Errorf("expected ONLY_A in OnlyInA, got %+v", r.OnlyInA)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := []Entry{{Key: "SHARED", Value: "x"}}
	b := []Entry{{Key: "SHARED", Value: "x"}, {Key: "ONLY_B", Value: "2"}}
	r := Compare(a, b)
	if len(r.OnlyInB) != 1 || r.OnlyInB[0].Key != "ONLY_B" {
		t.Errorf("expected ONLY_B in OnlyInB, got %+v", r.OnlyInB)
	}
}

func TestCompare_Mixed(t *testing.T) {
	a := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "old"}, {Key: "C", Value: "same"}}
	b := []Entry{{Key: "B", Value: "new"}, {Key: "C", Value: "same"}, {Key: "D", Value: "4"}}
	r := Compare(a, b)
	if len(r.OnlyInA) != 1 || r.OnlyInA[0].Key != "A" {
		t.Errorf("OnlyInA mismatch: %+v", r.OnlyInA)
	}
	if len(r.OnlyInB) != 1 || r.OnlyInB[0].Key != "D" {
		t.Errorf("OnlyInB mismatch: %+v", r.OnlyInB)
	}
	if len(r.Changed) != 1 || r.Changed[0].Key != "B" {
		t.Errorf("Changed mismatch: %+v", r.Changed)
	}
	if len(r.Identical) != 1 || r.Identical[0].Key != "C" {
		t.Errorf("Identical mismatch: %+v", r.Identical)
	}
}

func TestFormatCompareResult_NoDiff(t *testing.T) {
	r := CompareResult{Identical: []Entry{{Key: "X", Value: "1"}}}
	out := FormatCompareResult(r, "a", "b")
	if !strings.Contains(out, "Identical: 1") {
		t.Errorf("expected identical count in output, got: %s", out)
	}
}

func TestFormatCompareResult_Empty(t *testing.T) {
	r := CompareResult{}
	out := FormatCompareResult(r, "a", "b")
	if out != "No differences found.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
