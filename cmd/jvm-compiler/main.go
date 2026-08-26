package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
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
		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}
		expr, err := parse(scanner.Text())
		if err != nil {
			return err
		}
		classFile, err := compile(expr)
		if err != nil {
			return err
		}
		if _, err := w.Write(classFile); err != nil {
			return fmt.Errorf("error: write class file: %w", err)
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
	return expression{left, matches[2][0], right}, nil
}

func parseInt32(value string) (int32, error) {
	n, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error: invalid 32-bit signed integer %q", value)
	}
	return int32(n), nil
}

type constantPool struct {
	entries [][]byte
}

func (p *constantPool) utf8(value string) uint16 {
	var entry bytes.Buffer
	entry.WriteByte(1)
	u2(&entry, uint16(len(value)))
	entry.WriteString(value)
	return p.add(entry.Bytes())
}

func (p *constantPool) ref(tag byte, indexes ...uint16) uint16 {
	var entry bytes.Buffer
	entry.WriteByte(tag)
	for _, index := range indexes {
		u2(&entry, index)
	}
	return p.add(entry.Bytes())
}

func (p *constantPool) integer(value int32) uint16 {
	var entry bytes.Buffer
	entry.WriteByte(3)
	u4(&entry, uint32(value))
	return p.add(entry.Bytes())
}

func (p *constantPool) add(entry []byte) uint16 {
	p.entries = append(p.entries, entry)
	return uint16(len(p.entries))
}

func compile(expr expression) ([]byte, error) {
	pool := &constantPool{}
	thisClass := pool.ref(7, pool.utf8("com/nurkiewicz/PL0"))
	superClass := pool.ref(7, pool.utf8("java/lang/Object"))
	codeName := pool.utf8("Code")
	mainName := pool.utf8("main")
	mainDescriptor := pool.utf8("([Ljava/lang/String;)V")
	systemClass := pool.ref(7, pool.utf8("java/lang/System"))
	outType := pool.ref(12, pool.utf8("out"), pool.utf8("Ljava/io/PrintStream;"))
	systemOut := pool.ref(9, systemClass, outType)
	printStreamClass := pool.ref(7, pool.utf8("java/io/PrintStream"))
	printlnType := pool.ref(12, pool.utf8("println"), pool.utf8("(I)V"))
	println := pool.ref(10, printStreamClass, printlnType)
	sourceFileName := pool.utf8("SourceFile")
	sourceFile := pool.utf8("PL0.pl0")

	code := []byte{0xb2, byte(systemOut >> 8), byte(systemOut)}
	code = appendPush(code, expr.left, pool)
	code = appendPush(code, expr.right, pool)
	opcode, ok := map[byte]byte{'+': 0x60, '-': 0x64, '*': 0x68, '/': 0x6c}[expr.op]
	if !ok {
		return nil, fmt.Errorf("error: unsupported operator %q", expr.op)
	}
	code = append(code, opcode, 0xb6, byte(println>>8), byte(println), 0xb1)

	var class bytes.Buffer
	u4(&class, 0xcafebabe)
	u2(&class, 0)
	u2(&class, 52)
	u2(&class, uint16(len(pool.entries)+1))
	for _, entry := range pool.entries {
		class.Write(entry)
	}
	u2(&class, 0x0021)
	u2(&class, thisClass)
	u2(&class, superClass)
	u2(&class, 0)
	u2(&class, 0)
	u2(&class, 1)
	method(&class, mainName, mainDescriptor, codeName, 3, 1, code)
	u2(&class, 1)
	u2(&class, sourceFileName)
	u4(&class, 2)
	u2(&class, sourceFile)
	return class.Bytes(), nil
}

func appendPush(code []byte, value int32, pool *constantPool) []byte {
	switch {
	case value == -1:
		return append(code, 0x02)
	case value >= 0 && value <= 5:
		return append(code, byte(0x03+value))
	case value >= -128 && value <= 127:
		return append(code, 0x10, byte(value))
	case value >= -32768 && value <= 32767:
		return append(code, 0x11, byte(value>>8), byte(value))
	default:
		index := pool.integer(value)
		if index <= 0xff {
			return append(code, 0x12, byte(index))
		}
		return append(code, 0x13, byte(index>>8), byte(index))
	}
}

func method(w *bytes.Buffer, name, descriptor, codeName, maxStack, maxLocals uint16, code []byte) {
	u2(w, 0x0009)
	u2(w, name)
	u2(w, descriptor)
	u2(w, 1)
	u2(w, codeName)
	u4(w, uint32(12+len(code)))
	u2(w, maxStack)
	u2(w, maxLocals)
	u4(w, uint32(len(code)))
	w.Write(code)
	u2(w, 0)
	u2(w, 0)
}

func u2(w *bytes.Buffer, value uint16) {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], value)
	w.Write(b[:])
}

func u4(w *bytes.Buffer, value uint32) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], value)
	w.Write(b[:])
}
