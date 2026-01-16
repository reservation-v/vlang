package modfile

import (
	"fmt"
	"strings"
)

func ParseModulePath(data []byte) (string, error) {
	stringArr := strings.Split(string(data), "\n")

	for _, line := range stringArr {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "module ") {
			line = strings.TrimSpace(line[len("module "):])
			return line, nil
		}
	}
	return "", fmt.Errorf("module not found in go.mod")

}
