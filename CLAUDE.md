# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build ./...          # build
go test ./...           # run all tests
go test -run TestName   # run a single test
go vet ./...            # lint
go run main.go          # run (reads from stdin)
echo "3 + 4" | go run main.go
go run main.go < test.txt
```

## Architecture

This is a single-file Go CLI interpreter (`main.go`) with no external dependencies.

**Data flow:** `main()` → `run(io.Reader)` → `interpret(line string)` per non-blank line → prints result to stdout.

- `interpret` strips all whitespace, scans past the first number token (including an optional leading sign) to find the operator (`+`, `-`, `*`, `/`), then parses both sides as `float64`. Division by zero is an error.
- `run` is the loop: scans lines, skips blanks, calls `interpret`, and returns on the first error.
- `test.txt` contains sample input for manual testing.
