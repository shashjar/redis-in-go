package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// TODO: finish implementing Redis protocol
/**
 * Implementation of the Redis Serialization Protocol (RESP): https://redis.io/docs/latest/develop/reference/protocol-spec/
 */

const (
	SIMPLE_STRING   = "+"
	SIMPLE_ERROR    = "-"
	INTEGER         = ":"
	BULK_STRING     = "$"
	ARRAY           = "*"
	NULL            = "_"
	BOOLEAN         = "#"
	DOUBLE          = ","
	BIG_NUMBER      = "("
	BULK_ERROR      = "!"
	VERBATIM_STRING = "="
	MAP             = "%"
	SET             = "~"
	PUSH            = ">"
)

func toSimpleString(s string) string {
	return SIMPLE_STRING + s + "\r\n"
}

func toSimpleError(errorMessage string) string {
	return SIMPLE_ERROR + errorMessage + "\r\n"
}

func toInteger(num int) string {
	var sign string = ""
	if num < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%s%d\r\n", INTEGER, sign, num)
}

func toBulkString(s string) string {
	return fmt.Sprintf("%s%d\r\n%s\r\n", BULK_STRING, len(s), s)
}

func toNullBulkString() string {
	return BULK_STRING + "-1\r\n"
}

func toArray(a []string) string {
	arrayString := fmt.Sprintf("%s%d\r\n", ARRAY, len(a))
	for _, s := range a {
		arrayString += toBulkString(s)
	}
	return arrayString
}

func expectRESPDataType(b []byte, pos int, expectedType string) (int, error) {
	if pos >= len(b) {
		return pos, io.ErrUnexpectedEOF
	}

	if b[pos] != expectedType[0] {
		return pos, errors.New("unexpected RESP data type")
	}

	return pos + 1, nil
}

func expectBytes(b []byte, pos int, expected []byte) (int, error) {
	if pos+len(expected) > len(b) {
		return pos, io.ErrUnexpectedEOF
	}

	for i := 0; i < len(expected); i++ {
		if b[pos+i] != expected[i] {
			return pos, fmt.Errorf("expected byte %s, got %s", string(expected[i]), string(b[pos+i]))
		}
	}

	return pos + len(expected), nil
}

func parseInteger(b []byte, pos int) (int, int, error) {
	i := pos
	if i >= len(b) {
		return 0, i, io.ErrUnexpectedEOF
	}

	var n int
	for {
		if i >= len(b) || b[i] < '0' || b[i] > '9' {
			break
		}

		digit, err := strconv.Atoi(string(b[i]))
		if err != nil {
			return 0, i, err
		}
		n = 10*n + digit
		i++
	}

	i, err := expectBytes(b, i, []byte{'\r', '\n'})
	if err != nil {
		return 0, i, err
	}

	return n, i, nil
}

func parseBulkString(b []byte, pos int) (string, int, error) {
	i := pos
	if i >= len(b) {
		return "", i, io.ErrUnexpectedEOF
	}

	n, i, err := parseInteger(b, i)
	if err != nil {
		return "", i, err
	}

	if i+n > len(b) {
		return "", i, io.ErrUnexpectedEOF
	}

	s := string(b[i : i+n])
	i += n
	i, err = expectBytes(b, i, []byte{'\r', '\n'})
	if err != nil {
		return "", i, err
	}

	return s, i, nil
}

func parseCommand(b []byte) ([]string, error) {
	var commandComponents []string

	i := 0
	if i >= len(b) {
		return nil, io.ErrUnexpectedEOF
	}

	// Array length
	i, err := expectRESPDataType(b, i, ARRAY)
	if err != nil {
		return nil, err
	}
	n, i, err := parseInteger(b, i)
	if err != nil {
		return nil, err
	}

	// Array elements (command components are bulk strings)
	for j := 0; j < n; j++ {
		var commandComponent string

		i, err = expectRESPDataType(b, i, BULK_STRING)
		if err != nil {
			return nil, err
		}
		commandComponent, i, err = parseBulkString(b, i)
		if err != nil {
			return nil, err
		}
		commandComponents = append(commandComponents, commandComponent)
	}

	return commandComponents, nil
}
