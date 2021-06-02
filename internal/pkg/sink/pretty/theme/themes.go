package theme

// Default ...
var Default = prince

// Registry is the registry of theme names to their structs.
var Registry = map[string]Theme{}

// Names ...
var Names []string

func register(name string, theme Theme) {
	Names = append(Names, name)
	Registry[name] = theme
}

func style(fg string, bg string) columnStyle {
	return columnStyle{
		Fg: fg,
		Bg: bg,
	}
}
