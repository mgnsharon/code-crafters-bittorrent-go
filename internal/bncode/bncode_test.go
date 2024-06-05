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
		{"decode a list within a list", "l4:textli50e4:testei50ei51e5:helloe", []interface{}{"text", []interface{}{50, "test"}, 50, 51, "hello"}},
		{"decode a dictionary", "d3:foo3:bar5:helloi52ee", map[string]interface{}{"foo": "bar", "hello": 52}},
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

func TestBncode(t *testing.T) {
	tcs := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "test", "4:test"},
		{"int", 50, "i50e"},
		{"list", []interface{}{"test", 50, "hello"}, "l4:testi50e5:helloe"},
		{"dictionary", map[string]interface{}{"foo": "bar", "hello": 52}, "d3:foo3:bar5:helloi52ee"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := bncode.Bncode(tc.input)
			assertNoError(t, err)
			if got != tc.expected {
				t.Errorf("got %s, want %s", got, tc.expected)
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
