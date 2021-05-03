package manager

import (
	"couture/internal/pkg/model"
)

func VerboseDisplayOption(verbosity uint) interface{} {
	return baseOption{applier: func(m *managerOptions) {
		switch verbosity {
		case 0:
			m.level = model.LevelError
		case 1:
			m.level = model.LevelWarn
		case 2:
			m.level = model.LevelInfo
		case 3:
			m.level = model.LevelDebug
		default:
			m.level = model.LevelTrace
		}
	}}
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
		level      model.Level
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
