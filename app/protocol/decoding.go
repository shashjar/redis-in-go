package protocol

import (
	"io"
	"strconv"
)

func ParseCommand(b []byte) ([]string, error) {
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
