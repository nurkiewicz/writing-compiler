## Commands

```bash
go build ./...                      # build
go build -o compiler ./cmd/compiler # build compiler
make all                            # build compiler
go test ./...                       # run all tests
go vet ./...                        # lint
echo "3 + 4" | go run ./cmd/compiler | xxd
```

## Architecture

This project contains a compiler CLI executable with no external dependencies.

`cmd/compiler` parses one arithmetic expression and writes binary bytecode to stdout.

**Data flow:** stdin → compiler → binary bytecode.
