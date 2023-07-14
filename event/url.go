package event

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// SourceURL ...
type SourceURL url.URL

// String ...
func (u *SourceURL) String() string {
	u2 := url.URL(*u)
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
		if flag, err := strconv.ParseBool(s); err == nil {
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
// In the first case this results in a hostname of some and a path of /path/list. This action rewrites it into the
// proper second form.
func (u *SourceURL) Normalize() {
	if u.Host != "" {
		u.Path = "/" + u.Host + u.Path
		u.Host = ""
	}
}

// ShortForm ...
func (u *SourceURL) ShortForm() string {
	const tldComponentCount = 2
	var host = u.Host
	u.RawQuery = strings.TrimRight(u.RawQuery, "&")
	if net.ParseIP(host) == nil {
		hostParts := strings.Split(host, ".")
		if len(hostParts) > tldComponentCount {
			host = strings.Join(hostParts[0:len(hostParts)-tldComponentCount], ".")
		}
	}
	path := strings.Split(strings.TrimLeft(u.Path, "/"), "/")[0]
	if path == "" {
		if host == "" {
			return fmt.Sprintf("%s[%s]", u.Scheme, u.RawQuery)
		}
		return fmt.Sprintf("%s[%s?%s]", u.Scheme, host, u.RawQuery)
	}
	return fmt.Sprintf("%s[%s/%s]", u.Scheme, host, path)
}

// Hash ...
func (u *SourceURL) Hash() int {
	var sum int
	for _, v := range sha256.Sum256([]byte(u.String())) {
		sum += int(v)
	}
	return sum
}

// HashString ...
func (u *SourceURL) HashString() string {
	hasher := sha256.New()
	hasher.Write([]byte(u.String()))
	return strings.ReplaceAll(hex.EncodeToString(hasher.Sum(nil)), "-", "")
}
