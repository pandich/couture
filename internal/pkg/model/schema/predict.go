package schema

// Guess ..
func Guess(possibleJSON string, schemasToCheck ...Schema) *Schema {
	for _, schema := range schemasToCheck {
		if schema.CanHandle(possibleJSON) {
			return &schema
		}
	}
	return nil
}
