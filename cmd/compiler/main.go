package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const pushOpcode byte = 0x01

var header = []byte{'P', 'L', '/', '0', 0x00, 0x01}

type expression struct {
	left  int32
	op    byte
	right int32
}

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		expr, err := parse(line)
		if err != nil {
			return err
		}
		return writeProgram(w, expr)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return errors.New("error: empty expression")
}

const (
	ws       = `\s*`
	number   = `[+-]?\d+`
	operator = `[+\-*/]`
)

func c(pattern string) string {
	return `(` + pattern + `)`
}

var exprRegex = regexp.MustCompile(`^` + ws + c(number) + ws + c(operator) + ws + c(number) + ws + `$`)

func parse(line string) (expression, error) {
	matches := exprRegex.FindStringSubmatch(line)
	if matches == nil {
		if strings.TrimSpace(line) == "" {
			return expression{}, errors.New("error: empty expression")
		}
		return expression{}, fmt.Errorf("error: expected \"integer op integer\", got %q", line)
	}

	left, err := parseInt32(matches[1])
	if err != nil {
		return expression{}, err
	}
	right, err := parseInt32(matches[3])
	if err != nil {
		return expression{}, err
	}

	return expression{left: left, op: matches[2][0], right: right}, nil
}

func parseInt32(value string) (int32, error) {
	n, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error: invalid 32-bit signed integer %q", value)
	}
	return int32(n), nil
}

func writeProgram(w io.Writer, expr expression) error {
	program := make([]byte, 0, len(header)+11)
	program = append(program, header...)
	program = appendPush(program, expr.left)
	program = appendPush(program, expr.right)
	program = append(program, expr.op)

	if _, err := w.Write(program); err != nil {
		return fmt.Errorf("error: write program: %w", err)
	}
	return nil
}

func appendPush(program []byte, value int32) []byte {
	program = append(program, pushOpcode)
	return binary.BigEndian.AppendUint32(program, uint32(value))
}
