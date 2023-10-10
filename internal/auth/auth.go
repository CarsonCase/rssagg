package auth

import (
	"errors"
	"net/http"
	"strings"
)

// extracts api key from headers of http request
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("No auth info")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("Invalid auth info")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("Invalid auth info")
	}

	return vals[1], nil
}
