package main

import (
	"bytes"
	"testing"
)

func TestRunWritesProgram(t *testing.T) {
	var got bytes.Buffer
	if err := run(bytes.NewBufferString("40 + 2\n"), &got); err != nil {
		t.Fatalf("run() error: %v", err)
	}

	want := []byte{
		'P', 'L', '/', '0', 0x00, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x28,
		0x01, 0x00, 0x00, 0x00, 0x02,
		'+',
	}
	if !bytes.Equal(got.Bytes(), want) {
		t.Errorf("run() = % x, want % x", got.Bytes(), want)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  expression
	}{
		{"3 + 4", expression{3, '+', 4}},
		{"40+2", expression{40, '+', 2}},
		{"-5 - -3", expression{-5, '-', -3}},
		{"2 * 4", expression{2, '*', 4}},
		{"10 / 4", expression{10, '/', 4}},
	}
	for _, tt := range tests {
		got, err := parse(tt.input)
		if err != nil {
			t.Errorf("parse(%q) error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("parse(%q) = %+v, want %+v", tt.input, got, tt.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	tests := []string{
		"",
		"42",
		"1.5 + 2",
		"2147483648 + 1",
		"hello",
		"+ 2",
		"2 +",
	}
	for _, input := range tests {
		if _, err := parse(input); err == nil {
			t.Errorf("parse(%q) expected error, got nil", input)
		}
	}
}

func TestRunEmptyInput(t *testing.T) {
	var output bytes.Buffer
	if err := run(bytes.NewBuffer(nil), &output); err == nil {
		t.Fatal("run() expected error for empty input")
	}
}
