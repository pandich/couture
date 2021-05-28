package column

var allColumns = []column{
	newSourceColumn(),
	newTimestampColumn(),
	newApplicationColumn(),
	newThreadColumn(),
	newCallerColumn(),
	newLevelColumn(),
	newMessageColumn(),
}

func init() {
	// build the columnRegistry
	for _, col := range allColumns {
		registry[col.name()] = col
	}
}

// Names all available column names.
func Names() []string {
	var columnNames []string
	for _, col := range allColumns {
		columnNames = append(columnNames, col.name())
	}
	return columnNames
}

var registry = columnRegistry{}

type columnRegistry map[string]column
