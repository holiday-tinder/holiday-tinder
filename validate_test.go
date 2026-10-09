package main

import "testing"

func TestNewCodeIsValid(t *testing.T) {
	for i := 0; i < 200; i++ {
		c, err := newCode()
		if err != nil {
			t.Fatal(err)
		}
		if !validCode(c) {
			t.Fatalf("generated invalid code %q", c)
		}
	}
}

func TestValidCode(t *testing.T) {
	cases := map[string]bool{"ABC234": true, "abc234": false, "ABCO12": false, "ABC": false, "": false}
	for in, want := range cases {
		if got := validCode(in); got != want {
			t.Errorf("validCode(%q) = %v, want %v", in, got, want)
		}
	}
	if normalizeCode(" abc234 ") != "ABC234" {
		t.Error("normalizeCode should trim and uppercase")
	}
}

func TestCleanName(t *testing.T) {
	if n, ok := cleanName("  Sam "); !ok || n != "Sam" {
		t.Errorf("cleanName trimmed = %q %v", n, ok)
	}
	if _, ok := cleanName("   "); ok {
		t.Error("blank name accepted")
	}
	if _, ok := cleanName("abcdefghijklmnopqrstuvwxyzabcdefg"); ok {
		t.Error("overlong name accepted")
	}
}

func TestCleanCategories(t *testing.T) {
	out, ok := cleanCategories([]string{"Bar", "club", "bar"})
	if !ok || len(out) != 2 {
		t.Errorf("got %v %v", out, ok)
	}
	if _, ok := cleanCategories([]string{"casino"}); ok {
		t.Error("unknown category accepted")
	}
	if _, ok := cleanCategories(nil); ok {
		t.Error("empty categories accepted")
	}
}

func TestIsUUID(t *testing.T) {
	if !isUUID("3f2504e0-4f89-11d3-9a0c-0305e82c3301") {
		t.Error("valid uuid rejected")
	}
	if isUUID("3f2504e0-4f89-11d3-9a0c-0305e82c330z") || isUUID("nope") {
		t.Error("invalid uuid accepted")
	}
}
