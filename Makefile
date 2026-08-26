.PHONY: all compiler vm jvm-compiler clean

all: compiler vm jvm-compiler

compiler:
	go build -o compiler ./cmd/compiler

vm:
	go build -o vm ./cmd/vm

jvm-compiler:
	go build -o jvm-compiler ./cmd/jvm-compiler

clean:
	rm -f compiler vm jvm-compiler
