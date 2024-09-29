package protocol

import (
	"errors"
	"fmt"
	"io"
)

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
