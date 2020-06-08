package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"gopkg.in/validator.v2"
)

// MustReadFromEnv populates and validates the given object with values from env.
// It will panic if there's an error or the config is invalid
func MustReadFromEnv(cfg interface{}) {
	if err := env.Parse(cfg); err != nil {
		panic(fmt.Errorf("cannot parse env config: %w", err))
	}

	if errs := validator.Validate(cfg); errs != nil {
		panic(fmt.Errorf("invalid env config: %w", errs))
	}
}
