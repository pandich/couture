package model

import (
	"strings"
	"time"
)

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

// Stamp ...
func (t Timestamp) Stamp() string {
	return time.Time(t).Format(time.Stamp)
}
