package mongodbcluster

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

// mergeHelmValuesMap merges the newMap into baseMap.
func mergeHelmValuesMap(baseMap pulumi.Map, newMap map[string]string) {
	for key, val := range newMap {
		setNestedValue(baseMap, key, val)
	}
}

// helper function to navigate through a pulumi.Map and set a value based on a path.
// Note: This is a conceptual example; direct manipulation like this with pulumi.Map is not typical.
func setNestedValue(baseMap pulumi.Map, key string, value interface{}) {
	parts := strings.Split(key, ".")
	last := len(parts) - 1
	current := baseMap

	for i, part := range parts {
		// Convert current part to pulumi.Map if it's not already.
		if _, ok := current[part].(pulumi.Map); !ok && i < last {
			current[part] = pulumi.Map{}
		}

		if i == last {
			// Set the value at the last part.
			current[part] = pulumi.String(value.(string))
		} else {
			// Navigate deeper into the map.
			current = current[part].(pulumi.Map)
		}
	}
}
