package types

// Rule struct so that we get traefik config.
type Rule struct {
	Name         string   `yaml:"name"`
	Header       string   `yaml:"header"`
	Value        string   `yaml:"value"`
	Values       []string `yaml:"values"`
	HeaderPrefix string   `yaml:"headerPrefix"`
	Sep          string   `yaml:"sep"`
}
