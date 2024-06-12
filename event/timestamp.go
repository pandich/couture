package event

import (
	"strings"
	"time"
)

// HumanTimeFormat procides human-friendly times via humanize.Time.
const HumanTimeFormat = "human"

// Timestamp When the even occurred.
type Timestamp time.Time

// UnmarshalJSON ...
func (t *Timestamp) UnmarshalJSON(bytes []byte) error {
	tsString := strings.Trim(string(bytes), `"`)
	ts, err := time.Parse(time.RFC3339, tsString)
	if err != nil {
		return err
	}
	*t = Timestamp(ts)
	return nil
}
