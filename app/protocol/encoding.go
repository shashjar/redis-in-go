package protocol

import "fmt"

func ToSimpleString(s string) string {
	return SIMPLE_STRING + s + "\r\n"
}

func ToSimpleError(errorMessage string) string {
	return SIMPLE_ERROR + errorMessage + "\r\n"
}

func ToInteger(num int) string {
	var sign string = ""
	if num < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%s%d\r\n", INTEGER, sign, num)
}

func ToBulkString(s string) string {
	return fmt.Sprintf("%s%d\r\n%s\r\n", BULK_STRING, len(s), s)
}

func ToNullBulkString() string {
	return BULK_STRING + "-1\r\n"
}

func ToArray(a []string) string {
	arrayString := fmt.Sprintf("%s%d\r\n", ARRAY, len(a))
	for _, s := range a {
		arrayString += ToBulkString(s)
	}
	return arrayString
}
