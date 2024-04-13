// Package traefik_header_rename transforms header keys for Traefik
package traefik_header_rename

import (
	"context"
	"github.com/dsdeboer/traefik-header-rename/lib"
	"net/http"
	"regexp"
)

// New created a new HeaderRenamer plugin.
func New(_ context.Context, next http.Handler, config *lib.Config, name string) (http.Handler, error) {
	return &HeaderRenamer{
		rules: config.Rules,
		next:  next,
		name:  name,
	}, nil
}

// HeaderRenamer holds the necessary components of a Traefik plugin.
type HeaderRenamer struct {
	next  http.Handler
	rules []lib.Rule
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
					req.Header.Set(lib.GetValue(rule.Value, rule.HeaderPrefix, req), val)
				}
			}
		}
	}
	u.next.ServeHTTP(rw, req)
}
