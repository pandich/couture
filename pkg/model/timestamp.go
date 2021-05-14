package model

import (
	"strings"
	"time"
)

// Timestamp ...
type (
	// Timestamp When the even occurred.
	Timestamp time.Time
	// Stamp ...
	Stamp string
)

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
func (t Timestamp) Stamp() Stamp {
	return Stamp(time.Time(t).Format(time.Stamp))
}
