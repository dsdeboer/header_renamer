// Package traefik_header_rename transforms header keys for Traefik
package traefik_header_rename

import (
	"context"
	"net/http"
	"regexp"
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
			matched, err := regexp.MatchString(rule.OldHeader, headerName)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)

				return
			}
			if matched {
				req.Header.Del(headerName)
				for _, val := range headerValues {
					req.Header.Set(rule.NewHeader, val)
				}
			}
		}
	}
	u.next.ServeHTTP(rw, req)
}

// Config holds configuration to be passed to the plugin.
type Config struct {
	Rules []Rule
}

// Rule struct so that we get traefik config.
type Rule struct {
	OldHeader string `yaml:"old-header"`
	NewHeader string `yaml:"new-header"`
}
