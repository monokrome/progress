package progress

import (
	"github.com/monokrome/prefer.go"
)

// StorageOptions contains options describing where data should be stored
type StorageOptions struct {
	Backend string `yaml:"backend" json:"backend"`
	Options string `yaml:"options" json:"options"`
}

// Options stores user options
type Options struct {
	Storage StorageOptions `yaml:"storage" json:"storage"`
}

// NewOptions creates a new Options structure
func NewOptions(identifier string) (Options, *prefer.Configuration, error) {
	options := Options{}
	configuration, err := prefer.Load("progress", &options)

	if err != nil {
		return options, nil, err
	}

	return options, configuration, nil
}
