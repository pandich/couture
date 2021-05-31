package schema

// Guess ..
func Guess(s string, schemasToCheck ...Schema) *Schema {
	for _, schema := range schemasToCheck {
		if schema.CanHandle(s) {
			return &schema
		}
	}
	return nil
}
