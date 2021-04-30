package cli

import (
	"errors"
	"fmt"
	"gopkg.in/multierror.v1"
	"time"
)

const (
	//maxLineCount is the inclusive upper bound of coreCli.LineCount
	maxLineCount uint32 = 5_000
	//minPollInterval is the shortest allowed interval for polling sources.
	minPollInterval = 1 * time.Second
	//maxPollInterval is the longest allowed interval for polling sources.
	maxPollInterval = 10 * time.Minute
)

//coreValidator validates coreCli.
type coreValidator struct{}

func (c coreValidator) Validate() error {
	var violations []error
	if coreCli.LineCount > maxLineCount {
		violations = append(violations, fmt.Errorf("line count may not be greater than %d", maxLineCount))
	}
	if len(Sources()) == 0 {
		violations = append(violations, errors.New("at least one source must be specified"))
	}
	if len(Sinks()) == 0 {
		violations = append(violations, errors.New("at least one destination must be specified"))
	}
	if coreCli.PollInterval < minPollInterval || coreCli.PollInterval > maxPollInterval {
		violations = append(violations, fmt.Errorf("interval must be >= %v and <= %v", minPollInterval, maxPollInterval))
	}
	if len(violations) > 0 {
		return multierror.New(violations)
	}
	return nil
}
