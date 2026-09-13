// Package resp implements the Redis Serialization Protocol (RESP2).
//
// RESP2 types:
//
//	+OK\r\n           -> Simple String
//	-ERR ...\r\n     -> Error
//	:123\r\n         -> Integer
//	$5\r\nhello\r\n -> Bulk String
//	*2\r\n...       -> Array
//
// The protocol is line-oriented with \r\n terminators.
// Bulk strings and arrays are prefixed with a length.
// Null bulk strings are represented as $-1\r\n.
// Null arrays are represented as *-1\r\n.
package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Value is a parsed RESP value.
type Value struct {
	typ   Type
	null  bool    // true for null bulk strings and null arrays
	str   string  // SimpleString, Error, BulkString
	num   int64   // Integer
	array []Value // Array elements
}

// Type represents a RESP value type.
type Type byte

const (
	SimpleString Type = '+'
	Error        Type = '-'
	Integer      Type = ':'
	BulkString   Type = '$'
	Array        Type = '*'
)

// Constructors
func NewSimpleString(s string) Value { return Value{typ: SimpleString, str: s} }
func NewError(s string) Value        { return Value{typ: Error, str: s} }
func NewInteger(n int64) Value       { return Value{typ: Integer, num: n} }
func NewBulkString(s string) Value   { return Value{typ: BulkString, str: s} }
func NewNullBulkString() Value       { return Value{typ: BulkString, null: true} }
func NewArray(vs ...Value) Value     { return Value{typ: Array, array: vs} }
func NewNullArray() Value            { return Value{typ: Array, null: true} }

// Getters
func (v Value) Type() Type     { return v.typ }
func (v Value) String() string { return v.str }
func (v Value) Error() string  { return v.str }
func (v Value) Integer() int64 { return v.num }
func (v Value) Array() []Value { return v.array }
func (v Value) IsNull() bool   { return v.null }

// String returns a human-readable representation for debugging.
func (v Value) DebugString() string {
	switch v.typ {
	case SimpleString:
		return fmt.Sprintf("(+ %q)", v.str)
	case Error:
		return fmt.Sprintf("(- %q)", v.str)
	case Integer:
		return fmt.Sprintf("(: %d)", v.num)
	case BulkString:
		if v.IsNull() {
			return "($ nil)"
		}
		return fmt.Sprintf("($ %q)", v.str)
	case Array:
		if v.IsNull() {
			return "(* nil)"
		}
		return fmt.Sprintf("(* %v)", v.array)
	default:
		return fmt.Sprintf("(? %v)", v)
	}
}

// Reader reads RESP values from an io.Reader.
type Reader struct {
	reader *bufio.Reader
}

// NewReader creates a new RESP reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(r)}
}

// ReadValue reads and parses the next RESP value.
// It supports both RESP arrays and inline commands (space-separated text lines).
func (r *Reader) ReadValue() (Value, error) {
	// Read the type byte.
	b, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch Type(b) {
	case SimpleString:
		return r.readSimpleString()
	case Error:
		return r.readError()
	case Integer:
		return r.readInteger()
	case BulkString:
		return r.readBulkString()
	case Array:
		return r.readArray()
	default:
		// Inline command: put the byte back and read the whole line.
		if err := r.reader.UnreadByte(); err != nil {
			return Value{}, err
		}
		return r.readInline()
	}
}

// readLine reads until \r\n and returns the line without the terminator.
func (r *Reader) readLine() ([]byte, error) {
	line, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(line) < 2 || line[len(line)-2] != '\r' {
		return nil, fmt.Errorf("resp: invalid line ending")
	}
	return line[:len(line)-2], nil
}

func (r *Reader) readSimpleString() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	return NewSimpleString(string(line)), nil
}

func (r *Reader) readError() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	return NewError(string(line)), nil
}

func (r *Reader) readInteger() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	n, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return Value{}, fmt.Errorf("resp: invalid integer: %w", err)
	}
	return NewInteger(n), nil
}

func (r *Reader) readBulkString() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	length, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return Value{}, fmt.Errorf("resp: invalid bulk string length: %w", err)
	}
	if length < 0 {
		// Null bulk string.
		return NewNullBulkString(), nil
	}

	// Read exactly length bytes + \r\n.
	data := make([]byte, length+2)
	_, err = io.ReadFull(r.reader, data)
	if err != nil {
		return Value{}, fmt.Errorf("resp: bulk string read: %w", err)
	}
	if data[length] != '\r' || data[length+1] != '\n' {
		return Value{}, fmt.Errorf("resp: invalid bulk string terminator")
	}
	return NewBulkString(string(data[:length])), nil
}

func (r *Reader) readArray() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	count, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return Value{}, fmt.Errorf("resp: invalid array length: %w", err)
	}
	if count < 0 {
		// Null array.
		return NewNullArray(), nil
	}

	values := make([]Value, count)
	for i := int64(0); i < count; i++ {
		v, err := r.ReadValue()
		if err != nil {
			return Value{}, err
		}
		values[i] = v
	}
	return NewArray(values...), nil
}

// readInline reads a space-separated command line and converts it to a RESP array.
// Example: "GET foo\r\n" -> Array[BulkString("GET"), BulkString("foo")]
func (r *Reader) readInline() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	// Split by spaces. We don't handle quoted arguments for inline commands.
	parts := strings.Fields(string(line))
	if len(parts) == 0 {
		return NewArray(), nil
	}
	values := make([]Value, len(parts))
	for i, p := range parts {
		values[i] = NewBulkString(p)
	}
	return NewArray(values...), nil
}

// Writer writes RESP values to an io.Writer.
type Writer struct {
	writer io.Writer
}

// NewWriter creates a new RESP writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// WriteValue serializes a RESP value.
func (w *Writer) WriteValue(v Value) error {
	switch v.typ {
	case SimpleString:
		return w.writeSimpleString(v.str)
	case Error:
		return w.writeError(v.str)
	case Integer:
		return w.writeInteger(v.num)
	case BulkString:
		if v.IsNull() {
			_, err := fmt.Fprint(w.writer, "$-1\r\n")
			return err
		}
		return w.writeBulkString(v.str)
	case Array:
		return w.writeArray(v.array, v.null)
	default:
		return fmt.Errorf("resp: unknown type: %v", v.typ)
	}
}

func (w *Writer) writeSimpleString(s string) error {
	_, err := fmt.Fprintf(w.writer, "+%s\r\n", s)
	return err
}

func (w *Writer) writeError(s string) error {
	_, err := fmt.Fprintf(w.writer, "-%s\r\n", s)
	return err
}

func (w *Writer) writeInteger(n int64) error {
	_, err := fmt.Fprintf(w.writer, ":%d\r\n", n)
	return err
}

func (w *Writer) writeBulkString(s string) error {
	_, err := fmt.Fprintf(w.writer, "$%d\r\n%s\r\n", len(s), s)
	return err
}

func (w *Writer) writeArray(values []Value, isNull bool) error {
	if isNull {
		_, err := fmt.Fprint(w.writer, "*-1\r\n")
		return err
	}
	_, err := fmt.Fprintf(w.writer, "*%d\r\n", len(values))
	if err != nil {
		return err
	}
	for _, v := range values {
		if err := w.WriteValue(v); err != nil {
			return err
		}
	}
	return nil
}
