## Commands

```bash
go build ./...                      # build
go build -o compiler ./cmd/compiler # build compiler
go build -o vm ./cmd/vm             # build VM
make all                            # build compiler and VM
go test ./...                       # run all tests
go vet ./...                        # lint
echo "3 + 4" | go run ./cmd/compiler | go run ./cmd/vm
```

## Architecture

This project contains two Go CLI executables with no external dependencies:

- `cmd/compiler` parses one arithmetic expression and writes binary bytecode to stdout.
- `cmd/vm` reads bytecode from stdin, executes it, and prints the result to stdout.

**Data flow:** stdin → compiler → binary bytecode → VM → result.
