package kv

import (
	"fmt"
	"strings"
)

// Parse parses a key-value string in the format: key1='value';key2='value2';key3=0
// Returns a map of key-value pairs.
func Parse(input string) (map[string]string, error) {
	if input == "" {
		return make(map[string]string), nil
	}

	result := make(map[string]string)
	
	// Split by semicolon to get individual key=value pairs
	pairs := strings.Split(input, ";")
	
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		
		// Find the first '=' to split key and value
		idx := strings.Index(pair, "=")
		if idx == -1 {
			return nil, fmt.Errorf("invalid format: missing '=' in '%s'", pair)
		}
		
		key := strings.TrimSpace(pair[:idx])
		value := strings.TrimSpace(pair[idx+1:])
		
		if key == "" {
			return nil, fmt.Errorf("invalid format: empty key in '%s'", pair)
		}
		
		// Remove surrounding quotes (single or double)
		value = unquote(value)
		
		result[key] = value
	}
	
	return result, nil
}

// unquote removes surrounding single or double quotes from a string
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
