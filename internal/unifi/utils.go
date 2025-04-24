package unifi

import (
	"strings"
)

// FormatUrl formats a URL with the given parameters.
func FormatUrl(path string, params ...string) string {
	segments := strings.Split(path, "%s")
	for i, param := range params {
		if param != "" {
			segments[i] += param
		}
	}
	return strings.Join(segments, "")
}

func Find[T any](slice []T, predicate func(T) bool) *T {
	for i := range slice {
		if predicate(slice[i]) {
			return &slice[i]
		}
	}
	return nil
}

func RemoveWhere[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if !predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
