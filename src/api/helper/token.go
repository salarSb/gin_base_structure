package helper

import (
	"errors"
	"strings"
)

func ExtractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}
	fields := strings.Fields(authHeader)
	switch len(fields) {
	case 1:
		return fields[0], nil
	case 2:
		if strings.EqualFold(fields[0], "Bearer") {
			return fields[1], nil
		}
	}
	return "", errors.New("invalid authorization header format")
}
