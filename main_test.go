package main

import (
	"strings"
	"testing"
)

func TestInterpret(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"3 + 4", 7},
		{"40+2", 42},
		{"40 +2", 42},
		{"40+ 2", 42},
		{"  40  +  2  ", 42},
		{"0 + 0", 0},
		{"1000000 + 2000000", 3000000},
		{"-5 + 3", -2},
		{"5 + -3", 2},
		{"-5 + -3", -8},
		{"1.5 + 2.5", 4.0},
		{"0.1 + 0.2", 0.30000000000000004},
		{"-3.14 + 3.14", 0},
	}
	for _, tt := range tests {
		got, err := interpret(tt.input)
		if err != nil {
			t.Errorf("interpret(%q) error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("interpret(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestInterpretErrors(t *testing.T) {
	tests := []string{
		"",
		"42",
		"hello",
		"abc + def",
		"+ 2",
		"2 +",
		"10 - 3",
		"3 * 4",
		"10 / 4",
	}
	for _, input := range tests {
		_, err := interpret(input)
		if err == nil {
			t.Errorf("interpret(%q) expected error, got nil", input)
		}
	}
}

func TestRun(t *testing.T) {
	input := "3 + 4\n10+20\n"
	got, err := run(strings.NewReader(input))
	if err != nil {
		t.Fatalf("run() error: %v", err)
	}
	if got != 7 {
		t.Errorf("run() = %v, want 7", got)
	}
}

func TestRunError(t *testing.T) {
	input := "bad\n"
	_, err := run(strings.NewReader(input))
	if err == nil {
		t.Fatal("run() expected error for invalid input")
	}
}
