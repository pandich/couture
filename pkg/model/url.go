package model

import (
	"errors"
	"fmt"
	"github.com/araddon/dateparse"
	errors2 "github.com/pkg/errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	errBadSince = errors.New("bad value for since parameter")
)

// SourceURL ...
type SourceURL url.URL

// Since ...
func (u SourceURL) Since(key string) (*time.Time, error) {
	var since *time.Time
	var c = url.URL(u)
	for k, v := range c.Query() {
		if k != key {
			continue
		}
		if len(v) == 0 {
			return nil, errors2.WithMessagef(errBadSince, "may not be blank")
		}
		var err error
		var d time.Duration
		arg := v[0]
		d, err = time.ParseDuration(arg)
		if err == nil {
			t := time.Now().Add(-d)
			since = &t
		} else {
			var t time.Time
			t, err = dateparse.ParseAny(arg)
			if err == nil {
				since = &t
			}
		}
		if err != nil {
			return nil, errors2.WithMessagef(errBadSince, "could not parse - %s", arg)
		}
	}
	return since, nil
}

// String ...
func (u SourceURL) String() string {
	u2 := url.URL(u)
	return u2.String()
}

// QueryKey looks up a QueryKey parameter value.
func (u *SourceURL) QueryKey(key string) (string, bool) {
	c := url.URL(*u)
	for k, v := range c.Query() {
		if k == key {
			if len(v) > 0 {
				return v[0], true
			}
			return "", true
		}
	}
	return "", false
}

// QueryFlag looks up a QueryKey parameter value.
// If the key exists with an empty value, it is the equivalent of having set the value to true.
func (u *SourceURL) QueryFlag(key string) bool {
	var v string
	var ok bool
	if v, ok = u.QueryKey(key); ok {
		var s = strings.Trim(v, " ")
		if s == "" {
			return true
		}
		var err error
		var flag bool
		if flag, err = strconv.ParseBool(s); err == nil {
			return flag
		}
	}
	return false
}

// QueryInt64 returns the value of the QueryKey parameter a pointer to an int64. If the parameter
// is not set, or is set to empty, a nil pointer with no error is returned.
func (u *SourceURL) QueryInt64(key string) (*int64, error) {
	s, ok := u.QueryKey(key)
	if !ok {
		return nil, nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// Normalize fixes situations where scheme://some/path/list is given rather than scheme:///some/path/list.
// In the first case this results in a hostname of some and a path of /path/list. This method rewrites it into the
// proper second form.
func (u *SourceURL) Normalize() {
	if u.Host != "" {
		u.Path = "/" + u.Host + u.Path
		u.Host = ""
	}
}

// ShortForm ...
func (u SourceURL) ShortForm() string {
	const tldComponentCount = 2
	var host = u.Host
	if net.ParseIP(host) == nil {
		hostParts := strings.Split(host, ".")
		if len(hostParts) > tldComponentCount {
			host = strings.Join(hostParts[0:len(hostParts)-tldComponentCount], ".")
		}
	}
	path := strings.Split(strings.TrimLeft(u.Path, "/"), "/")[0]
	return fmt.Sprintf(":/%s/%s", host, path)
}