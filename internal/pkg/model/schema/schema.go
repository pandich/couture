package schema

// Logstash ...
var Logstash = newSchema(mapping{
	"@timestamp":           TimestampCol,
	"level":                LevelCol,
	"message":              MessageCol,
	"application":          ApplicationCol,
	"method":               MethodCol,
	"line_number":          LineCol,
	"thread_name":          ThreadCol,
	"class":                ClassCol,
	"exception.stacktrace": ExceptionCol,
})

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

type mapping map[Field]column

// Schema ...
type Schema struct {
	Mapping mapping
	Fields  []string
}

func newSchema(mapping mapping) Schema {
	var fields []string
	for f := range mapping {
		fields = append(fields, string(f))
	}
	return Schema{
		Mapping: mapping,
		Fields:  fields,
	}
}
