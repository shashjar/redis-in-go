package protocol

import "fmt"

// Serializes to the RESP Simple String type
func ToSimpleString(s string) string {
	return SIMPLE_STRING + s + "\r\n"
}

// Serializes to the RESP Simple Error type
func ToSimpleError(errorMessage string) string {
	return SIMPLE_ERROR + errorMessage + "\r\n"
}

// Serializes to the RESP Integer type
func ToInteger(num int) string {
	var sign string = ""
	if num < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%s%d\r\n", INTEGER, sign, num)
}

// Serializes to the RESP Bulk String type
func ToBulkString(s string) string {
	return fmt.Sprintf("%s%d\r\n%s\r\n", BULK_STRING, len(s), s)
}

// Serializes to the RESP Null Bulk String type
func ToNullBulkString() string {
	return BULK_STRING + "-1\r\n"
}

// Serializes to the RESP Array type
func ToArray(a []string) string {
	arrayString := fmt.Sprintf("%s%d\r\n", ARRAY, len(a))
	for _, s := range a {
		arrayString += ToBulkString(s)
	}
	return arrayString
}
