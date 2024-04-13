package lib

import (
	"net/http"
	"strings"
)

// GetValue checks if prefix exists
// the given prefix is present, and then proceeds to read
// the existing header (after stripping the prefix) to return as value.
func GetValue(ruleValue, valueIsHeaderPrefix string, req *http.Request) string {
	actualValue := ruleValue
	if valueIsHeaderPrefix != "" && strings.HasPrefix(ruleValue, valueIsHeaderPrefix) {
		header := strings.TrimPrefix(ruleValue, valueIsHeaderPrefix)
		// If the resulting value after removing the prefix is empty (value was only prefix),
		// we return the actual value, which is the prefix itself.
		// This is because doing a req.Header.Get("") would not fly well.
		if header != "" {
			actualValue = req.Header.Get(header)
		}
	}

	return actualValue
}
