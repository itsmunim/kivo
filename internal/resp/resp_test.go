package resp

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Serialization Tests ---

func TestWriteSimpleString(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewSimpleString("OK"))
	require.NoError(t, err)
	assert.Equal(t, "+OK\r\n", buf.String())
}

func TestWriteError(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewError("ERR unknown command 'foo'"))
	require.NoError(t, err)
	assert.Equal(t, "-ERR unknown command 'foo'\r\n", buf.String())
}

func TestWriteInteger(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewInteger(42))
	require.NoError(t, err)
	assert.Equal(t, ":42\r\n", buf.String())
}

func TestWriteIntegerNegative(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewInteger(-999))
	require.NoError(t, err)
	assert.Equal(t, ":-999\r\n", buf.String())
}

func TestWriteBulkString(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewBulkString("hello"))
	require.NoError(t, err)
	assert.Equal(t, "$5\r\nhello\r\n", buf.String())
}

func TestWriteBulkStringEmpty(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewBulkString(""))
	require.NoError(t, err)
	assert.Equal(t, "$0\r\n\r\n", buf.String())
}

func TestWriteNullBulkString(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewNullBulkString())
	require.NoError(t, err)
	assert.Equal(t, "$-1\r\n", buf.String())
}

func TestWriteArray(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewArray(
		NewBulkString("GET"),
		NewBulkString("key"),
	))
	require.NoError(t, err)
	assert.Equal(t, "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n", buf.String())
}

func TestWriteEmptyArray(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewArray())
	require.NoError(t, err)
	assert.Equal(t, "*0\r\n", buf.String())
}

func TestWriteNullArray(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewNullArray())
	require.NoError(t, err)
	assert.Equal(t, "*-1\r\n", buf.String())
}

func TestWriteNestedArray(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteValue(NewArray(
		NewInteger(1),
		NewArray(
			NewBulkString("a"),
			NewBulkString("b"),
		),
		NewInteger(3),
	))
	require.NoError(t, err)
	assert.Equal(t, "*3\r\n:1\r\n*2\r\n$1\r\na\r\n$1\r\nb\r\n:3\r\n", buf.String())
}

// --- Parsing Tests ---

func TestReadSimpleString(t *testing.T) {
	r := NewReader(strings.NewReader("+OK\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, SimpleString, v.Type())
	assert.Equal(t, "OK", v.String())
}

func TestReadError(t *testing.T) {
	r := NewReader(strings.NewReader("-ERR unknown command 'foo'\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Error, v.Type())
	assert.Equal(t, "ERR unknown command 'foo'", v.Error())
}

func TestReadInteger(t *testing.T) {
	r := NewReader(strings.NewReader(":42\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Integer, v.Type())
	assert.Equal(t, int64(42), v.Integer())
}

func TestReadIntegerNegative(t *testing.T) {
	r := NewReader(strings.NewReader(":-999\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, int64(-999), v.Integer())
}

func TestReadBulkString(t *testing.T) {
	r := NewReader(strings.NewReader("$5\r\nhello\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, BulkString, v.Type())
	assert.Equal(t, "hello", v.String())
	assert.False(t, v.IsNull())
}

func TestReadBulkStringEmpty(t *testing.T) {
	r := NewReader(strings.NewReader("$0\r\n\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, BulkString, v.Type())
	assert.Equal(t, "", v.String())
	assert.False(t, v.IsNull()) // Empty string is not null.
}

func TestReadNullBulkString(t *testing.T) {
	r := NewReader(strings.NewReader("$-1\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, BulkString, v.Type())
	assert.True(t, v.IsNull())
}

func TestReadArray(t *testing.T) {
	r := NewReader(strings.NewReader("*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v.Type())
	arr := v.Array()
	require.Len(t, arr, 2)
	assert.Equal(t, "GET", arr[0].String())
	assert.Equal(t, "key", arr[1].String())
}

func TestReadEmptyArray(t *testing.T) {
	r := NewReader(strings.NewReader("*0\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v.Type())
	assert.Empty(t, v.Array())
}

func TestReadNullArray(t *testing.T) {
	r := NewReader(strings.NewReader("*-1\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v.Type())
	assert.True(t, v.IsNull())
}

func TestReadNestedArray(t *testing.T) {
	r := NewReader(strings.NewReader("*3\r\n:1\r\n*2\r\n$1\r\na\r\n$1\r\nb\r\n:3\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v.Type())
	arr := v.Array()
	require.Len(t, arr, 3)
	assert.Equal(t, int64(1), arr[0].Integer())
	assert.Equal(t, Array, arr[1].Type())
	assert.Equal(t, int64(3), arr[2].Integer())
}

// --- Round-trip Tests ---

func TestRoundTripSimpleString(t *testing.T) {
	original := NewSimpleString("hello world")
	assertRoundTrip(t, original)
}

func TestRoundTripError(t *testing.T) {
	original := NewError("ERR something went wrong")
	assertRoundTrip(t, original)
}

func TestRoundTripInteger(t *testing.T) {
	original := NewInteger(-12345)
	assertRoundTrip(t, original)
}

func TestRoundTripBulkString(t *testing.T) {
	original := NewBulkString("this is a longer string with spaces and symbols !@#")
	assertRoundTrip(t, original)
}

func TestRoundTripNullBulkString(t *testing.T) {
	original := NewNullBulkString()
	assertRoundTrip(t, original)
}

func TestRoundTripArray(t *testing.T) {
	original := NewArray(
		NewBulkString("SET"),
		NewBulkString("mykey"),
		NewBulkString("myvalue"),
		NewInteger(42),
	)
	assertRoundTrip(t, original)
}

func TestRoundTripComplexArray(t *testing.T) {
	original := NewArray(
		NewBulkString("LPUSH"),
		NewBulkString("mylist"),
		NewArray(
			NewBulkString("nested"),
			NewInteger(1),
		),
		NewNullBulkString(),
	)
	assertRoundTrip(t, original)
}

func assertRoundTrip(t *testing.T, original Value) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	require.NoError(t, w.WriteValue(original))

	r := NewReader(&buf)
	parsed, err := r.ReadValue()
	require.NoError(t, err)
	assertEqualValue(t, original, parsed)
}

func assertEqualValue(t *testing.T, expected, actual Value) {
	t.Helper()
	assert.Equal(t, expected.Type(), actual.Type(), "type mismatch")
	switch expected.Type() {
	case SimpleString, Error, BulkString:
		assert.Equal(t, expected.String(), actual.String(), "string mismatch")
	case Integer:
		assert.Equal(t, expected.Integer(), actual.Integer(), "integer mismatch")
	case Array:
		if expected.IsNull() {
			assert.True(t, actual.IsNull(), "expected null array")
			return
		}
		expArr := expected.Array()
		actArr := actual.Array()
		require.Len(t, actArr, len(expArr), "array length mismatch")
		for i := range expArr {
			assertEqualValue(t, expArr[i], actArr[i])
		}
	}
}

// --- Pipelining Tests ---

func TestReadPipelinedValues(t *testing.T) {
	input := "+OK\r\n:42\r\n$5\r\nhello\r\n*-1\r\n"
	r := NewReader(strings.NewReader(input))

	v1, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, SimpleString, v1.Type())
	assert.Equal(t, "OK", v1.String())

	v2, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, int64(42), v2.Integer())

	v3, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, "hello", v3.String())

	v4, err := r.ReadValue()
	require.NoError(t, err)
	assert.True(t, v4.IsNull())
}

// --- Error Cases ---

func TestReadUnknownTypeByte(t *testing.T) {
	// With inline command support, '?' is treated as start of an inline command.
	// This test verifies that non-RESP type bytes are handled as inline.
	r := NewReader(strings.NewReader("?unknown\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v.Type())
	arr := v.Array()
	require.Len(t, arr, 1)
	assert.Equal(t, "?unknown", arr[0].String())
}

func TestReadInvalidInteger(t *testing.T) {
	r := NewReader(strings.NewReader(":notanumber\r\n"))
	_, err := r.ReadValue()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid integer")
}

func TestReadInvalidBulkStringLength(t *testing.T) {
	r := NewReader(strings.NewReader("$abc\r\n"))
	_, err := r.ReadValue()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid bulk string length")
}

func TestReadBulkStringTooShort(t *testing.T) {
	r := NewReader(strings.NewReader("$10\r\nshort\r\n"))
	_, err := r.ReadValue()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bulk string read")
}

func TestReadBulkStringBadTerminator(t *testing.T) {
	r := NewReader(strings.NewReader("$5\r\nhello\n\n"))
	_, err := r.ReadValue()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid bulk string terminator")
}

func TestReadInvalidLineEnding(t *testing.T) {
	r := NewReader(strings.NewReader("+OK\n"))
	_, err := r.ReadValue()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid line ending")
}

func TestReadEOF(t *testing.T) {
	r := NewReader(strings.NewReader(""))
	_, err := r.ReadValue()
	assert.ErrorIs(t, err, io.EOF)
}

func TestReadPartialArray(t *testing.T) {
	r := NewReader(strings.NewReader("*2\r\n$3\r\nGET\r\n"))
	_, err := r.ReadValue()
	require.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

// --- Benchmarks ---

func BenchmarkWriteBulkString(b *testing.B) {
	w := NewWriter(io.Discard)
	v := NewBulkString("hello world this is a test value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = w.WriteValue(v)
	}
}

func BenchmarkReadBulkString(b *testing.B) {
	data := []byte("$32\r\nhello world this is a test value\r\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(bytes.NewReader(data))
		_, _ = r.ReadValue()
	}
}

func BenchmarkWriteArray(b *testing.B) {
	w := NewWriter(io.Discard)
	v := NewArray(
		NewBulkString("SET"),
		NewBulkString("mykey"),
		NewBulkString("myvalue"),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = w.WriteValue(v)
	}
}

func BenchmarkReadArray(b *testing.B) {
	data := []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(bytes.NewReader(data))
		_, _ = r.ReadValue()
	}
}

// --- Inline Command Tests ---

func TestReadInlineSimple(t *testing.T) {
	r := NewReader(strings.NewReader("PING\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v.Type())
	arr := v.Array()
	require.Len(t, arr, 1)
	assert.Equal(t, "PING", arr[0].String())
}

func TestReadInlineWithArgs(t *testing.T) {
	r := NewReader(strings.NewReader("GET mykey\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v.Type())
	arr := v.Array()
	require.Len(t, arr, 2)
	assert.Equal(t, "GET", arr[0].String())
	assert.Equal(t, "mykey", arr[1].String())
}

func TestReadInlineMultipleArgs(t *testing.T) {
	r := NewReader(strings.NewReader("SET mykey myvalue\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v.Type())
	arr := v.Array()
	require.Len(t, arr, 3)
	assert.Equal(t, "SET", arr[0].String())
	assert.Equal(t, "mykey", arr[1].String())
	assert.Equal(t, "myvalue", arr[2].String())
}

func TestReadInlineEmpty(t *testing.T) {
	r := NewReader(strings.NewReader("\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v.Type())
	assert.Empty(t, v.Array())
}

func TestReadInlineMixedWithRESP(t *testing.T) {
	// Inline PING followed by RESP array GET.
	input := "PING\r\n*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"
	r := NewReader(strings.NewReader(input))

	v1, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v1.Type())
	assert.Equal(t, "PING", v1.Array()[0].String())

	v2, err := r.ReadValue()
	require.NoError(t, err)
	assert.Equal(t, Array, v2.Type())
	assert.Equal(t, "GET", v2.Array()[0].String())
}

func TestReadInlineTelnetStyle(t *testing.T) {
	// Simulate telnet-style commands with multiple spaces.
	r := NewReader(strings.NewReader("SET   key   value\r\n"))
	v, err := r.ReadValue()
	require.NoError(t, err)
	arr := v.Array()
	require.Len(t, arr, 3)
	assert.Equal(t, "SET", arr[0].String())
	assert.Equal(t, "key", arr[1].String())
	assert.Equal(t, "value", arr[2].String())
}
