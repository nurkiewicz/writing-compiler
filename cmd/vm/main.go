package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	headerSize = 6
	pushOpcode = 0x01
)

var expectedHeader = []byte{'P', 'L', '/', '0', 0x00, 0x01}

func main() {
	result, err := execute(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, result)
}

func execute(r io.Reader) (int32, error) {
	program, err := io.ReadAll(r)
	if err != nil {
		return 0, fmt.Errorf("error: read program: %w", err)
	}
	if len(program) < headerSize {
		return 0, errors.New("error: missing or incomplete program header")
	}
	if string(program[:4]) != string(expectedHeader[:4]) {
		return 0, errors.New("error: unrecognized magic value")
	}
	if program[4] != expectedHeader[4] || program[5] != expectedHeader[5] {
		return 0, fmt.Errorf("error: unsupported bytecode version %d.%d", program[4], program[5])
	}

	var stack []int32
	for ip := headerSize; ip < len(program); {
		opcode := program[ip]
		ip++

		if opcode == pushOpcode {
			if len(program)-ip < 4 {
				return 0, errors.New("error: incomplete PUSH operand")
			}
			stack = append(stack, int32(binary.BigEndian.Uint32(program[ip:ip+4])))
			ip += 4
			continue
		}

		if opcode != '+' && opcode != '-' && opcode != '*' && opcode != '/' {
			return 0, fmt.Errorf("error: unknown opcode 0x%02x", opcode)
		}
		if len(stack) < 2 {
			return 0, fmt.Errorf("error: stack underflow for %q", opcode)
		}

		right := stack[len(stack)-1]
		left := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		var result int32
		switch opcode {
		case '+':
			result = left + right
		case '-':
			result = left - right
		case '*':
			result = left * right
		case '/':
			if right == 0 {
				return 0, errors.New("error: division by zero")
			}
			result = left / right
		}
		stack = append(stack, result)
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("error: expected one final stack value, got %d", len(stack))
	}
	return stack[0], nil
}
