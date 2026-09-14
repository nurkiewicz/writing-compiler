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

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

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
		module, err := compile(expr)
		if err != nil {
			return err
		}
		if _, err := io.WriteString(w, module.String()); err != nil {
			return fmt.Errorf("error: write LLVM IR: %w", err)
		}
		return nil
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return errors.New("error: empty expression")
}

var exprRegex = regexp.MustCompile(`^\s*([+-]?\d+)\s*([+\-*/])\s*([+-]?\d+)\s*$`)

func parse(line string) (expression, error) {
	matches := exprRegex.FindStringSubmatch(line)
	if matches == nil {
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

func parseInt32(input string) (int32, error) {
	n, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error: invalid 32-bit signed integer %q", input)
	}
	return int32(n), nil
}

func compile(expr expression) (*ir.Module, error) {
	module := ir.NewModule()

	format := module.NewGlobalDef("format", constant.NewCharArrayFromString("%d\n\x00"))
	format.Immutable = true

	printf := module.NewFunc("printf", types.I32, ir.NewParam("format", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true

	main := module.NewFunc("main", types.I32)
	entry := main.NewBlock("")
	left := constant.NewInt(types.I32, int64(expr.left))
	right := constant.NewInt(types.I32, int64(expr.right))

	var result value.Value
	switch expr.op {
	case '+':
		result = entry.NewAdd(left, right)
	case '-':
		result = entry.NewSub(left, right)
	case '*':
		result = entry.NewMul(left, right)
	case '/':
		result = entry.NewSDiv(left, right)
	default:
		return nil, fmt.Errorf("error: unsupported operator %q", expr.op)
	}

	formatPtr := entry.NewGetElementPtr(
		format.ContentType,
		format,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	entry.NewCall(printf, formatPtr, result)
	entry.NewRet(constant.NewInt(types.I32, 0))

	return module, nil
}
