package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWritesLLVMIR(t *testing.T) {
	var output bytes.Buffer
	if err := run(strings.NewReader("40 + 2\n"), &output); err != nil {
		t.Fatalf("run() error: %v", err)
	}

	ir := output.String()
	for _, want := range []string{
		`@format = constant [4 x i8] c"%d\0A\00"`,
		`declare i32 @printf(i8* %format, ...)`,
		`define i32 @main()`,
		`add i32 40, 2`,
		`call i32 (i8*, ...) @printf`,
		`ret i32 0`,
	} {
		if !strings.Contains(ir, want) {
			t.Errorf("LLVM IR does not contain %q:\n%s", want, ir)
		}
	}
}

func TestCompileArithmeticOperators(t *testing.T) {
	tests := []struct {
		expr        expression
		instruction string
	}{
		{expression{7, '+', 3}, "add i32 7, 3"},
		{expression{7, '-', 3}, "sub i32 7, 3"},
		{expression{7, '*', 3}, "mul i32 7, 3"},
		{expression{7, '/', 3}, "sdiv i32 7, 3"},
	}

	for _, test := range tests {
		module, err := compile(test.expr)
		if err != nil {
			t.Fatalf("compile(%+v) error: %v", test.expr, err)
		}
		if got := module.String(); !strings.Contains(got, test.instruction) {
			t.Errorf("compile(%+v) does not contain %q:\n%s", test.expr, test.instruction, got)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for _, input := range []string{"", "42", "1.5 + 2", "2147483648 + 1", "hello", "+ 2", "2 +"} {
		if _, err := parse(input); err == nil {
			t.Errorf("parse(%q) expected error, got nil", input)
		}
	}
}

func TestGeneratedIRCompilesAndRuns(t *testing.T) {
	clang, err := exec.LookPath("clang")
	if err != nil {
		t.Skip("clang is not installed")
	}

	var llvmIR bytes.Buffer
	if err := run(strings.NewReader("-20 * 2"), &llvmIR); err != nil {
		t.Fatalf("run() error: %v", err)
	}

	dir := t.TempDir()
	irPath := filepath.Join(dir, "expression.ll")
	executablePath := filepath.Join(dir, "expression")
	if err := os.WriteFile(irPath, llvmIR.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command(clang, irPath, "-o", executablePath).CombinedOutput(); err != nil {
		t.Fatalf("clang failed: %v\n%s", err, output)
	}

	output, err := exec.Command(executablePath).CombinedOutput()
	if err != nil {
		t.Fatalf("generated executable failed: %v\n%s", err, output)
	}
	if got := strings.TrimSpace(string(output)); got != "-40" {
		t.Errorf("executable output = %q, want %q", got, "-40")
	}
}
