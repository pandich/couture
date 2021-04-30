package manager

func VerboseDisplayOption() interface{} {
	return baseOption{f: func(m *managerOptions) error {
		m.verbose = true
		return nil
	}}
}
func QuietDisplayOption() interface{} {
	return baseOption{f: func(m *managerOptions) error {
		m.quiet = true
		return nil
	}}
}
func ClearScreenDisplayOption() interface{} {
	return baseOption{f: func(m *managerOptions) error {
		m.clearScreen = true
		return nil
	}}
}
func ShowPrefixDisplayOption() interface{} {
	return baseOption{f: func(m *managerOptions) error {
		m.showPrefix = true
		return nil
	}}
}
func ShortNamesDisplayOption() interface{} {
	return baseOption{f: func(m *managerOptions) error {
		m.shortNames = true
		return nil
	}}
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
		Apply(manager *managerOptions) error
	}

	baseOption struct {
		f func(*managerOptions) error
	}
)

func (b baseOption) Apply(mgr *managerOptions) error {
	return b.f(mgr)
}
