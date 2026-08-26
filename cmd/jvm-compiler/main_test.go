package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWritesJava8ClassWithJVMOpcodes(t *testing.T) {
	var output bytes.Buffer
	if err := run(strings.NewReader("2 + 3\n"), &output); err != nil {
		t.Fatalf("run() error: %v", err)
	}

	classFile := output.Bytes()
	if got := binary.BigEndian.Uint32(classFile[0:4]); got != 0xcafebabe {
		t.Fatalf("magic = %#x, want 0xcafebabe", got)
	}
	if got := binary.BigEndian.Uint16(classFile[6:8]); got != 52 {
		t.Errorf("major version = %d, want 52 (Java 8)", got)
	}

	// getstatic System.out, iconst_2, iconst_3, iadd, invokevirtual println, return
	wantCode := []byte{0xb2, 0x00, 0x0d, 0x05, 0x06, 0x60, 0xb6, 0x00, 0x13, 0xb1}
	if !bytes.Contains(classFile, wantCode) {
		t.Errorf("class file does not contain expected addition bytecode: % x", wantCode)
	}
	if !bytes.Contains(classFile, []byte("PL0.pl0")) {
		t.Error("class file does not contain PL0.pl0 source file name")
	}
}

func TestGeneratedClassRunsOnJava(t *testing.T) {
	java, err := exec.LookPath("java")
	if err != nil {
		t.Skip("java runtime is not installed")
	}
	if err := exec.Command(java, "-version").Run(); err != nil {
		t.Skip("java runtime is not usable")
	}

	var classFile bytes.Buffer
	if err := run(strings.NewReader("40 + 2"), &classFile); err != nil {
		t.Fatalf("run() error: %v", err)
	}

	root := t.TempDir()
	classDir := filepath.Join(root, "com", "nurkiewicz")
	if err := os.MkdirAll(classDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(classDir, "PL0.class"), classFile.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	command := exec.Command(java, "-cp", root, "com.nurkiewicz.PL0")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("java failed: %v\n%s", err, output)
	}
	if got := strings.TrimSpace(string(output)); got != "42" {
		t.Errorf("java output = %q, want %q", got, "42")
	}
}

func TestAppendPushUsesAppropriateOpcodes(t *testing.T) {
	tests := []struct {
		value int32
		want  []byte
	}{
		{-1, []byte{0x02}},
		{5, []byte{0x08}},
		{127, []byte{0x10, 0x7f}},
		{128, []byte{0x11, 0x00, 0x80}},
		{-32768, []byte{0x11, 0x80, 0x00}},
	}
	for _, test := range tests {
		pool := &constantPool{}
		if got := appendPush(nil, test.value, pool); !bytes.Equal(got, test.want) {
			t.Errorf("appendPush(%d) = % x, want % x", test.value, got, test.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for _, input := range []string{"", "42", "1.5 + 2", "2147483648 + 1", "hello"} {
		if _, err := parse(input); err == nil {
			t.Errorf("parse(%q) expected error, got nil", input)
		}
	}
}
