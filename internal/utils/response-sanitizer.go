package utils

import (
	"encoding/json"
)

func Sanitize(data []uint8) []uint8 {
	var tempMap map[string]interface{}
	err := json.Unmarshal(data, &tempMap)

	if err == nil {
		RemoveUserId(tempMap)
	}

	modifiedData, _ := json.Marshal(tempMap)
	return modifiedData
}

func RemoveUserId(data interface{}) {
	switch d := data.(type) {
	case map[string]interface{}:
		for key, value := range d {
			if key == "user_id" {
				delete(d, key)
			} else {
				RemoveUserId(value)
			}
		}
	case []interface{}:
		for _, value := range d {
			RemoveUserId(value)
		}
	}
}
