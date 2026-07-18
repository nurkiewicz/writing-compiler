.PHONY: all compiler vm clean

all: compiler vm

compiler:
	go build -o compiler ./cmd/compiler

vm:
	go build -o vm ./cmd/vm

clean:
	rm -f compiler vm
