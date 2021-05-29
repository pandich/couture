package schema

const (
	// Timestamp ...
	Timestamp = outputField("timestamp")
	// Level ...
	Level = outputField("level")
	// Message ...
	Message = outputField("message")
	// Application ...
	Application = outputField("application")
	// Method ...
	Method = outputField("method")
	// Line ...
	Line = outputField("line")
	// Thread ...
	Thread = outputField("thread")
	// Class ...
	Class = outputField("class")
	// Exception ...
	Exception = outputField("exception")
)

// InputField ...
type InputField string

type outputField string

// Mapping ...
type Mapping map[InputField]outputField

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
