package bncode_test

import (
	"encoding/json"
	"testing"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bncode"
)

func TestBencode(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected any
	}{
		{"decode string", "4:testwithsomeextratext", "test"},
		{"decode int stand alone", "i52e", 52},
		{"decode int extra string val", "i52e4:text", 52},
		{"decode a list", "l4:texti50ei51e5:helloe", []interface{}{"text", 50, 51, "hello"}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := bncode.Decode(tc.input)
			assertNoError(t, err)
			jsonOutput, _ := json.Marshal(got)
			expectedOutput, _ := json.Marshal(tc.expected)
			if string(jsonOutput) != string(expectedOutput) {
				t.Errorf("got %s, want %s", jsonOutput, expectedOutput)
			}
		})
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one: %v", got)
	}
}
