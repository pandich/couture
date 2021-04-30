package cli

import (
	"errors"
	"fmt"
	"gopkg.in/multierror.v1"
	"time"
)

const (
	//maxLineCount is the inclusive upper bound of coreCli.LineCount
	maxLineCount   uint32 = 5_000
	minPollCadence        = 1 * time.Second
	maxPollCadence        = 10 * time.Minute
)

//coreValidator validates coreCli.
type coreValidator struct{}

func (c coreValidator) Validate() error {
	var ea []error
	if coreCli.LineCount > maxLineCount {
		ea = append(ea, fmt.Errorf("line count may not be greater than %d", maxLineCount))
	}
	if len(Sources()) == 0 {
		ea = append(ea, errors.New("at least one source must be specified"))
	}
	if len(Sinks()) == 0 {
		ea = append(ea, errors.New("at least one destination must be specified"))
	}
	if coreCli.PollCadence < minPollCadence || coreCli.PollCadence > maxPollCadence {
		ea = append(ea, fmt.Errorf("interval must be >= %v and <= %v", minPollCadence, maxPollCadence))
	}
	if len(ea) > 0 {
		return multierror.New(ea)
	}
	return nil
}
