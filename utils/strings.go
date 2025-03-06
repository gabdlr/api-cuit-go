package utils

import "bytes"

func ReverseStringWithBuffer(input string) string {
	var buffer bytes.Buffer
	length := len(input) - 1
	for i := length; i >= 0; i-- {
		buffer.WriteByte(input[i])
	}
	return buffer.String()
}
