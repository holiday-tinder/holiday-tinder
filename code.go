package main

import (
	"crypto/rand"
	"strings"
)

// Unambiguous characters only (no 0/O, 1/I); 32 symbols divide 256 evenly.
const codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const codeLength = 6

func newCode() (string, error) {
	b := make([]byte, codeLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = codeAlphabet[int(b[i])%len(codeAlphabet)]
	}
	return string(b), nil
}

func validCode(s string) bool {
	if len(s) != codeLength {
		return false
	}
	for _, c := range s {
		if !strings.ContainsRune(codeAlphabet, c) {
			return false
		}
	}
	return true
}

func normalizeCode(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
