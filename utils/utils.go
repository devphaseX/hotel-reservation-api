package utils

type FailedResp struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// CloneMap creates a shallow copy of the original map
func CloneMap(original map[string]any) map[string]any {
	// Create a new map with the same type as the original
	clone := make(map[string]any, len(original))

	// Copy each key-value pair from the original to the new map
	for key, value := range original {
		clone[key] = value
	}

	return clone
}
