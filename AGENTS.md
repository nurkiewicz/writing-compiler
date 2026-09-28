## Commands

```bash
go build ./...                      # build
go build -o bin/compiler ./cmd/compiler # build compiler
go build -o bin/vm ./cmd/vm             # build VM
go build -o bin/jvm-compiler ./cmd/jvm-compiler # build JVM compiler
go build -o bin/llvm-compiler ./cmd/llvm-compiler # build LLVM compiler
make all                            # build all executables
go test ./...                       # run all tests
go vet ./...                        # lint
echo "3 + 4" | go run ./cmd/compiler | go run ./cmd/vm
```

## Architecture

This project contains four Go CLI executables:

- `cmd/compiler` parses one arithmetic expression and writes binary bytecode to stdout.
- `cmd/vm` reads bytecode from stdin, executes it, and prints the result to stdout.
- `cmd/jvm-compiler` parses the same expression and writes a Java 8
  `com.nurkiewicz.PL0` class file to stdout.
- `cmd/llvm-compiler` parses the same expression and writes textual LLVM IR
  containing a native `main` function to stdout.

**Data flow:** stdin → compiler → binary bytecode → VM → result.
