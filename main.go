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

const numberPattern = `\s*([+-]?(?:\d+\.?\d*|\.\d+))\s*`

var exprRegex = regexp.MustCompile(`^` + numberPattern + `([+])` + numberPattern + `$`)

func interpret(line string) (float64, error) {
	matches := exprRegex.FindStringSubmatch(line)
	if matches == nil {
		if strings.TrimSpace(line) == "" {
			return 0, errors.New("error: empty expression")
		}
		return 0, fmt.Errorf("error: expected \"number + number\", got %q", line)
	}

	left, _, right := matches[1], matches[2], matches[3]

	a, err := strconv.ParseFloat(left, 64)
	if err != nil {
		return 0, fmt.Errorf("error: invalid number %q", left)
	}
	b, err := strconv.ParseFloat(right, 64)
	if err != nil {
		return 0, fmt.Errorf("error: invalid number %q", right)
	}

	return a + b, nil
}
