package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	res, err := run(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(res)
}

func run(r io.Reader) (float64, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		return interpret(line)
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return 0, errors.New("error: empty expression")
}

const (
	ws       = `\s*`
	number   = `[+-]?(?:\d+\.?\d*|\.\d+)`
	operator = `[+\-*/]`
)

func c(pattern string) string {
	return `(` + pattern + `)`
}

var exprRegex = regexp.MustCompile(`^` + ws + c(number) + ws + c(operator) + ws + c(number) + ws + `$`)

func interpret(line string) (float64, error) {
	matches := exprRegex.FindStringSubmatch(line)
	if matches == nil {
		if strings.TrimSpace(line) == "" {
			return 0, errors.New("error: empty expression")
		}
		return 0, fmt.Errorf("error: expected \"number op number\", got %q", line)
	}

	left, opStr, right := matches[1], matches[2], matches[3]
	op := opStr[0]

	a, err := strconv.ParseFloat(left, 64)
	if err != nil {
		return 0, fmt.Errorf("error: invalid number %q", left)
	}
	b, err := strconv.ParseFloat(right, 64)
	if err != nil {
		return 0, fmt.Errorf("error: invalid number %q", right)
	}

	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, errors.New("error: division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("error: unknown operator %q", string(op))
	}
}
