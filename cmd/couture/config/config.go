package config

// config management to go here

type (
	SourceDefinition struct {
	}

	SinkDefinition struct {
	}

	Defaults struct {
	}

	Config struct {
		Sources  map[string]SourceDefinition
		Sinks    map[string]SinkDefinition
		Defaults Defaults
	}
)
