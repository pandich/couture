package schema

const (
	// TimestampCol ...
	TimestampCol = column("timestamp")
	// LevelCol ...
	LevelCol = column("level")
	// MessageCol ...
	MessageCol = column("message")
	// ApplicationCol ...
	ApplicationCol = column("application")
	// MethodCol ...
	MethodCol = column("method")
	// LineCol ...
	LineCol = column("line")
	// ThreadCol ...
	ThreadCol = column("thread")
	// ClassCol ...
	ClassCol = column("class")
	// ExceptionCol ...
	ExceptionCol = column("exception")
)

// Field ...
type Field string

type column string

// Mapping ...
type Mapping map[Field]column

// Schema ...
type Schema struct {
	Mapping Mapping
	Fields  []string
}

// NewSchema ...
func NewSchema(mapping Mapping) Schema {
	var fields []string
	for f := range mapping {
		fields = append(fields, string(f))
	}
	return Schema{
		Mapping: mapping,
		Fields:  fields,
	}
}
