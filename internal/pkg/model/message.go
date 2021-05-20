package model

import (
	"fmt"
	"regexp"
	"sort"
)

// TODO cleaner message/stacktrace approach

// Message ...
type (
	// Message a message.
	Message string

	// HighlightedMessage a highlighted section of a Message.
	HighlightedMessage string
	// UnhighlightedMessage ...
	UnhighlightedMessage string

	// Exception an exception.
	Exception struct {
		// StackTrace the full text of the stack trace.
		StackTrace StackTrace `json:"stack_trace"`
	}
	// StackTrace a stack trace.
	StackTrace Message
	// HighlightedStackTrace a highlighted section of a StackTrace.
	HighlightedStackTrace string
	// UnhighlightedStackTrace ...
	UnhighlightedStackTrace string
	// highlightMark ...

	highlightMark struct{ start, end int }
	// highlightMarks a collection of highlightMark elements.
	highlightMarks []highlightMark
)

// matches determines if an event matches the filters criteria.
func (msg Message) matches(include []regexp.Regexp, exclude []regexp.Regexp) (highlightMarks, bool) {
	var marks highlightMarks

	// process the includes returning true on the first match
	for _, filter := range include {
		for _, pos := range filter.FindAllStringIndex(string(msg), 100) {
			marks = append(marks, highlightMark{start: pos[0], end: pos[1]})
		}
	}
	// if we made it this far and have include filters, none of them matched, so we return false
	if len(include) > 0 && len(marks) == 0 {
		return nil, false
	}

	// process the excludes returning false on the first match
	for _, filter := range exclude {
		if filter.MatchString(string(msg)) {
			return nil, false
		}
	}
	// return true
	return marks, true
}

// String ...
func (msg Message) String() string {
	return string(msg)
}

// HighlightedMessage ...
func (msg Message) highlighted(
	allMarks highlightMarks,
	highlighted func(Message) interface{},
	unhighlighted func(message Message) interface{},
) []interface{} {
	if len(allMarks) == 0 {
		return []interface{}{unhighlighted(msg)}
	}

	var fields []interface{}
	var pos = 0
	marks := allMarks.merged()
	for _, mark := range marks {
		if mark.start > pos {
			fields = append(fields, unhighlighted(msg[pos:mark.start-1]))
		}
		pos = mark.end + 1
		fields = append(fields, highlighted(msg[mark.start:mark.end]))
	}
	if len(marks) > 0 {
		if marks[len(marks)-1].end < len(msg)-1 {
			fields = append(fields, unhighlighted(msg[marks[len(marks)-1].end+1:]))
		}
	}
	return fields
}

// merged merges all overlapping regions into contiguous ones.
// It is based upon https://stackoverflow.com/questions/55201821/merging-overlapping-intervals-using-double-for-loop
func (marks highlightMarks) merged() []highlightMark {
	m := append([]highlightMark(nil), marks...)
	if len(m) <= 1 {
		return m
	}

	sort.Slice(m,
		func(i, j int) bool {
			if m[i].start < m[j].start {
				return true
			}
			if m[i].start == m[j].start && m[i].end < m[j].end {
				return true
			}
			return false
		},
	)

	j := 0
	for i := 1; i < len(m); i++ {
		if m[j].end >= m[i].start {
			if m[j].end < m[i].end {
				m[j].end = m[i].end
			}
		} else {
			j++
			m[j] = m[i]
		}
	}
	return append([]highlightMark(nil), m[:j+1]...)
}

// NewException ...
func NewException(err error) *Exception {
	if err == nil {
		return nil
	}
	return &Exception{StackTrace: StackTrace(fmt.Sprint(err))}
}
