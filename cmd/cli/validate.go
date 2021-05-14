package cli

import (
	"errors"
	errors2 "github.com/pkg/errors"
	"gopkg.in/multierror.v1"
)

// cliValidator validates cli.
type cliValidator struct{}

// Validate provides validation for cli.
func (v cliValidator) Validate() error {
	var violations []error
	sources, err := configuredSources()
	if err != nil {
		violations = append(violations, errors2.Wrap(err, "sources could not be determined"))
	} else if len(sources) == 0 {
		violations = append(violations, errors.New("at least one source must be specified"))
	}
	if len(violations) > 0 {
		return multierror.New(violations)
	}
	return nil
}
