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
	indexes map[string]uint16
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
	if p.indexes == nil {
		p.indexes = make(map[string]uint16)
	}
	key := string(entry)
	if index, ok := p.indexes[key]; ok {
		return index
	}
	p.entries = append(p.entries, entry)
	index := uint16(len(p.entries))
	p.indexes[key] = index
	return index
}

func compile(expr expression) ([]byte, error) {
	pool := &constantPool{}
	thisClass := pool.ref(7, pool.utf8("com/nurkiewicz/PL0")) // CONSTANT_Class
	superClass := pool.ref(7, pool.utf8("java/lang/Object"))  // CONSTANT_Class
	codeName := pool.utf8("Code")
	mainName := pool.utf8("main")
	mainDescriptor := pool.utf8("([Ljava/lang/String;)V")
	systemClass := pool.ref(7, pool.utf8("java/lang/System"))                     // CONSTANT_Class
	outType := pool.ref(12, pool.utf8("out"), pool.utf8("Ljava/io/PrintStream;")) // CONSTANT_NameAndType
	systemOut := pool.ref(9, systemClass, outType)                                // CONSTANT_Fieldref
	printStreamClass := pool.ref(7, pool.utf8("java/io/PrintStream"))             // CONSTANT_Class
	printlnType := pool.ref(12, pool.utf8("println"), pool.utf8("(I)V"))          // CONSTANT_NameAndType
	println := pool.ref(10, printStreamClass, printlnType)                        // CONSTANT_Methodref
	sourceFileName := pool.utf8("SourceFile")
	sourceFile := pool.utf8("PL0.pl0")

	code := []byte{opcodeGetstatic, byte(systemOut >> 8), byte(systemOut)}
	code = appendPush(code, expr.left, pool)
	code = appendPush(code, expr.right, pool)
	opcode, ok := map[byte]byte{
		'+': opcodeIadd,
		'-': opcodeIsub,
		'*': opcodeImul,
		'/': opcodeIdiv,
	}[expr.op]
	if !ok {
		return nil, fmt.Errorf("error: unsupported operator %q", expr.op)
	}
	code = append(code, opcode, opcodeInvokevirtual, byte(println>>8), byte(println), opcodeReturn)

	var class bytes.Buffer
	u4(&class, 0xcafebabe)                  // magic
	u2(&class, 0)                           // minor_version
	u2(&class, 52)                          // major_version: Java 8
	u2(&class, uint16(len(pool.entries)+1)) // constant_pool_count
	for _, entry := range pool.entries {
		class.Write(entry)
	}
	u2(&class, 0x0021)                                             // public, super
	u2(&class, thisClass)                                          // this_class
	u2(&class, superClass)                                         // super_class
	u2(&class, 0)                                                  // interfaces_count
	u2(&class, 0)                                                  // fields_count
	u2(&class, 1)                                                  // methods_count
	method(&class, mainName, mainDescriptor, codeName, 3, 1, code) // max_stack, max_locals
	u2(&class, 1)                                                  // attributes_count
	u2(&class, sourceFileName)                                     // attribute_name_index
	u4(&class, 2)                                                  // attribute_length
	u2(&class, sourceFile)                                         // sourcefile_index
	return class.Bytes(), nil
}

func appendPush(code []byte, value int32, pool *constantPool) []byte {
	switch {
	case value == -1:
		return append(code, opcodeIconstM1)
	case value >= 0 && value <= 5:
		return append(code, byte(opcodeIconst0+value))
	case value >= -128 && value <= 127:
		return append(code, opcodeBipush, byte(value))
	case value >= -32768 && value <= 32767:
		return append(code, opcodeSipush, byte(value>>8), byte(value))
	default:
		index := pool.integer(value)
		if index <= 0xff {
			return append(code, opcodeLdc, byte(index))
		}
		return append(code, opcodeLdcW, byte(index>>8), byte(index))
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
