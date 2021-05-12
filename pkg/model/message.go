package model

import (
	"regexp"
	"sort"
)

// Message a message.
type Message string

// matches determines if an event matches the filters criteria.
func (msg Message) matches(include []*regexp.Regexp, exclude []*regexp.Regexp) (highlightMarks, bool) {
	var marks highlightMarks

	// process the includes returning true on the first match
	for _, filter := range include {
		for _, pos := range filter.FindAllStringIndex(string(msg), 100) {
			marks = append(marks, highlightMark{Start: pos[0], End: pos[1]})
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

// Highlighted ...
func (msg Message) highlighted(allMarks highlightMarks) []interface{} {
	var fields []interface{}
	var pos = 0
	marks := allMarks.merged()
	for _, mark := range marks {
		if mark.Start > pos {
			fields = append(fields, Unhighlighted(msg[pos:mark.Start-1]))
		}
		pos = mark.End + 1
		fields = append(fields, msg.highlight(mark))
	}
	if marks[len(marks)-1].End < len(msg)-1 {
		fields = append(fields, Unhighlighted(msg[marks[len(marks)-1].End+1:]))
	}
	return fields
}

// highlight ...
func (msg Message) highlight(mark highlightMark) Highlighted {
	return Highlighted(msg[mark.Start:mark.End])
}

// Highlighted a highlighted section of a Message.
type Highlighted string

// Unhighlighted ...
type Unhighlighted string

// highlightMark ...
type highlightMark struct {
	Start int
	End   int
}

// highlightMarks a collection of highlightMark elements.
type highlightMarks []highlightMark

// merged merges all overlapping regions into contiguous ones.
// It is based upon https://stackoverflow.com/questions/55201821/merging-overlapping-intervals-using-double-for-loop
func (marks highlightMarks) merged() []highlightMark {
	m := append([]highlightMark(nil), marks...)
	if len(m) <= 1 {
		return m
	}

	sort.Slice(m,
		func(i, j int) bool {
			if m[i].Start < m[j].Start {
				return true
			}
			if m[i].Start == m[j].Start && m[i].End < m[j].End {
				return true
			}
			return false
		},
	)

	j := 0
	for i := 1; i < len(m); i++ {
		if m[j].End >= m[i].Start {
			if m[j].End < m[i].End {
				m[j].End = m[i].End
			}
		} else {
			j++
			m[j] = m[i]
		}
	}
	return append([]highlightMark(nil), m[:j+1]...)
}
