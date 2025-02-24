package config

import (
	"bytes"
	"strings"
)

func FindLineAndCharacter(data []byte, offset int) (num, pos int) {
	lines := bytes.Split(data, []byte{'\n'})
	lineNumber := 1
	characterPosition := offset

	for _, line := range lines {
		if len(line)+1 < characterPosition {
			lineNumber++
			characterPosition -= len(line) + 1
		} else {
			break
		}
	}

	return lineNumber, characterPosition
}

func GetErrorContext(data []byte, offset int) string {
	start := offset - 20
	end := offset + 20

	if start < 0 {
		start = 0
	}

	if end > len(data) {
		end = len(data)
	}

	return strings.TrimSpace(string(data[start:end]))
}
