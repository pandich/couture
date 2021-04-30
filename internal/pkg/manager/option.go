package manager

func VerboseDisplayOption() interface{} {
	return baseOption{applier: func(m *managerOptions) { m.verbose = true }}
}
func QuietDisplayOption() interface{} {
	return baseOption{applier: func(m *managerOptions) { m.quiet = true }}
}
func ClearScreenDisplayOption() interface{} {
	return baseOption{applier: func(m *managerOptions) { m.clearScreen = true }}
}
func ShowPrefixDisplayOption() interface{} {
	return baseOption{applier: func(m *managerOptions) { m.showPrefix = true }}
}
func ShortNamesDisplayOption() interface{} {
	return baseOption{applier: func(m *managerOptions) { m.shortNames = true }}
}

type (
	managerOptions struct {
		verbose     bool
		quiet       bool
		clearScreen bool
		shortNames  bool
		showPrefix  bool
	}

	Option interface {
		Apply(manager *managerOptions)
	}

	baseOption struct {
		applier func(*managerOptions)
	}
)

func (b baseOption) Apply(mgrOptions *managerOptions) {
	b.applier(mgrOptions)
}
