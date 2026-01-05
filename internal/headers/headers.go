package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headerEndIndex := bytes.Index(data, []byte("\r\n"))
	if headerEndIndex == -1 {
		return 0, false, nil
	}
	if headerEndIndex == 0 {
		return 2, true, nil
	}
	fieldLineData := data[:headerEndIndex]
	separatorIndex := bytes.Index(fieldLineData, []byte(":"))
	if separatorIndex == -1 {
		return 0, false, fmt.Errorf("Field line does not contain separator.")
	}
	if separatorIndex == 0 {
		return 0, false, fmt.Errorf("Field name is empty.")
	}
	if fieldLineData[separatorIndex-1] == ' ' {
		return 0, false, fmt.Errorf("Field name must not end with whitespace.")
	}
	fieldName := string(bytes.Trim(fieldLineData[:separatorIndex], " "))
	fieldValue := string(bytes.Trim(fieldLineData[separatorIndex+1:], " "))
	h[fieldName] = fieldValue
	return len(fieldLineData) + 2, false, nil
}
