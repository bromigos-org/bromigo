package utils

import (
	"encoding/json"
	"log"
)

// logJSON is a helper function to log structs as JSON for debugging
func LogJSON(prefix string, v interface{}) {
	bytes, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Printf("Error marshalling %s to JSON: %v\n", prefix, err)
		return
	}
	log.Printf("%s JSON: %s\n", prefix, string(bytes))
}
