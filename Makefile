.PHONY: all clean

all: compiler vm jvm-compiler llvm-compiler

compiler: cmd/compiler/main.go go.mod go.sum
	go build -o $@ ./cmd/compiler

vm: cmd/vm/main.go go.mod go.sum
	go build -o $@ ./cmd/vm

jvm-compiler: cmd/jvm-compiler/main.go cmd/jvm-compiler/opcodes.go go.mod go.sum
	go build -o $@ ./cmd/jvm-compiler

llvm-compiler: cmd/llvm-compiler/main.go go.mod go.sum
	go build -o $@ ./cmd/llvm-compiler

clean:
	rm -f compiler vm jvm-compiler llvm-compiler
