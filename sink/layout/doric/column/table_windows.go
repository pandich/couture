//go:build windows

package column

// TODO windows dynamic column width support https://stackoverflow.com/a/10857339/2154219
func (table *Table) autoUpdateColumnWidths() {}
