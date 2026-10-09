package main

import (
	"strings"
	"unicode/utf8"
)

const maxNameLength = 30

var validCategories = map[string]bool{"restaurant": true, "bar": true, "club": true}

func cleanName(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" || utf8.RuneCountInString(s) > maxNameLength {
		return "", false
	}
	return s, true
}

func cleanCategories(in []string) ([]string, bool) {
	seen := map[string]bool{}
	out := []string{}
	for _, c := range in {
		c = strings.ToLower(strings.TrimSpace(c))
		if !validCategories[c] {
			return nil, false
		}
		if !seen[c] {
			seen[c] = true
			out = append(out, c)
		}
	}
	return out, len(out) > 0
}

func isUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
				return false
			}
		}
	}
	return true
}
