// Package traefik_header_rename transforms header keys for Traefik
package traefik_header_rename

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

// New created a new HeaderRenamer plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &HeaderRenamer{
		rules: config.Rules,
		next:  next,
		name:  name,
	}, nil
}

// CreateConfig populates the Config data object.
func CreateConfig() *Config {
	return &Config{
		Rules: []Rule{},
	}
}

// HeaderRenamer holds the necessary components of a Traefik plugin.
type HeaderRenamer struct {
	next  http.Handler
	rules []Rule
	name  string
}

func (u *HeaderRenamer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, rule := range u.rules {
		for headerName, headerValues := range req.Header {
			matched, err := regexp.MatchString(rule.Header, headerName)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)

				return
			}
			if matched {
				req.Header.Del(headerName)
				for _, val := range headerValues {
					req.Header.Set(GetValue(rule.Value, rule.HeaderPrefix, req), val)
				}
			}
		}
	}
	u.next.ServeHTTP(rw, req)
}

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

// Config holds configuration to be passed to the plugin.
type Config struct {
	Rules []Rule
}

// Rule struct so that we get traefik config.
type Rule struct {
	Name         string   `yaml:"name"`
	Header       string   `yaml:"header"`
	Value        string   `yaml:"value"`
	Values       []string `yaml:"values"`
	HeaderPrefix string   `yaml:"headerPrefix"`
	Sep          string   `yaml:"sep"`
}
