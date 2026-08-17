package main

import (
	"bytes"
	"testing"
)

func TestExecute(t *testing.T) {
	program := []byte{
		'P', 'L', '/', '0', 0x00, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x28,
		0x01, 0x00, 0x00, 0x00, 0x02,
		'+',
	}

	got, err := execute(bytes.NewReader(program))
	if err != nil {
		t.Fatalf("execute() error: %v", err)
	}
	if got != 42 {
		t.Errorf("execute() = %d, want 42", got)
	}
}

func TestExecuteErrors(t *testing.T) {
	tests := [][]byte{
		{},
		{'N', 'O', 'P', 'E', 0x00, 0x01},
		{'P', 'L', '/', '0', 0x00, 0x02},
		{'P', 'L', '/', '0', 0x00, 0x01, 0xff},
		{'P', 'L', '/', '0', 0x00, 0x01, '+'},
		{'P', 'L', '/', '0', 0x00, 0x01, 0x01},
	}
	for _, program := range tests {
		if _, err := execute(bytes.NewReader(program)); err == nil {
			t.Errorf("execute(% x) expected error, got nil", program)
		}
	}
}

func TestExecuteDivisionByZero(t *testing.T) {
	program := []byte{
		'P', 'L', '/', '0', 0x00, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00,
		'/',
	}
	if _, err := execute(bytes.NewReader(program)); err == nil {
		t.Fatal("execute() expected division by zero error")
	}
}
