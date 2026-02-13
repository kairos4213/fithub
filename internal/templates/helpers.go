package templates

import (
	"encoding/json"

	"github.com/kairos4213/fithub/internal/database"
)

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

func dateCompletedValue(workout database.Workout) string {
	if workout.DateCompleted.Valid {
		return workout.DateCompleted.Time.Format("2006-01-02")
	}
	return ""
}
