package templates

import "encoding/json"

// jsonVals safely marshals a map to a JSON string for use in hx-vals attributes.
// Using json.Marshal prevents injection when values contain special characters
// like double quotes, which would break hand-built JSON via fmt.Sprintf.
func jsonVals(data map[string]any) string {
	b, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(b)
}
