package modfile

import (
	"fmt"
	"strings"
)

func parseDirective(data []byte, key string) (string, error) {
	stringArr := strings.Split(string(data), "\n")

	for _, raw := range stringArr {
		line := strings.TrimSpace(raw)

		if idx := strings.Index(line, "//"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] != key {
			continue
		}
		if len(fields) != 2 {
			return "", fmt.Errorf("%s directive malformed", key)
		}
		if fields[1] == "" {
			return "", fmt.Errorf("%s directive has empty value", key)
		}
		return fields[1], nil
	}

	return "", fmt.Errorf("%s not found in go.mod", key)
}

func ParseModulePath(data []byte) (string, error) {
	return parseDirective(data, "module")
}

func ParseGoVersion(data []byte) (string, error) {
	return parseDirective(data, "go")
}
