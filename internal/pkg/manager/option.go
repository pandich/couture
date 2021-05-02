package manager

func VerboseDisplayOption() interface{} {
	return baseOption{applier: func(m *managerOptions) { m.verbose = true }}
}
func QuietDisplayOption() interface{} {
	return baseOption{applier: func(m *managerOptions) { m.quiet = true }}
}
func ShowPrefixDisplayOption() interface{} {
	return baseOption{applier: func(m *managerOptions) { m.showPrefix = true }}
}
func ShortNamesDisplayOption() interface{} {
	return baseOption{applier: func(m *managerOptions) { m.shortNames = true }}
}

type (
	//managerOptions
	managerOptions struct {
		verbose    bool
		quiet      bool
		shortNames bool // TODO move into sink.Options
		showPrefix bool // TODO move into sink.Options
	}

	Option interface {
		Apply(manager *managerOptions)
	}

	baseOption struct {
		applier func(*managerOptions)
	}
)

func (opt baseOption) Apply(mgrOptions *managerOptions) {
	opt.applier(mgrOptions)
}
