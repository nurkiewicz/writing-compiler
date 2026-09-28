BIN_DIR := bin
BINARIES := compiler vm jvm-compiler llvm-compiler
TARGETS := $(addprefix $(BIN_DIR)/,$(BINARIES))

.PHONY: all clean $(BINARIES)

all: $(TARGETS)

$(BINARIES): %: $(BIN_DIR)/%

$(BIN_DIR):
	mkdir -p $@

$(BIN_DIR)/compiler: cmd/compiler/main.go go.mod go.sum | $(BIN_DIR)
	go build -o $@ ./cmd/compiler

$(BIN_DIR)/vm: cmd/vm/main.go go.mod go.sum | $(BIN_DIR)
	go build -o $@ ./cmd/vm

$(BIN_DIR)/jvm-compiler: cmd/jvm-compiler/main.go cmd/jvm-compiler/opcodes.go go.mod go.sum | $(BIN_DIR)
	go build -o $@ ./cmd/jvm-compiler

$(BIN_DIR)/llvm-compiler: cmd/llvm-compiler/main.go go.mod go.sum | $(BIN_DIR)
	go build -o $@ ./cmd/llvm-compiler

clean:
	rm -f $(TARGETS)
	rmdir $(BIN_DIR) 2>/dev/null || true
