.PHONY: all compiler clean

all: compiler

compiler:
	go build -o compiler ./cmd/compiler

clean:
	rm -f compiler
