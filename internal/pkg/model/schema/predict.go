package schema

// Predict ..
func Predict(s string, schemas ...Schema) *Schema {
	for _, schema := range schemas {
		if schema.Test(s) {
			return &schema
		}
	}
	return nil
}
