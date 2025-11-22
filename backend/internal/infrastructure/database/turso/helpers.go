package database

import "encoding/json"

// parseStringArray JSON配列をパースするヘルパー関数
func parseStringArray(jsonStr string) []string {
	if jsonStr == "" {
		return []string{}
	}
	var arr []string
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		return []string{}
	}
	return arr
}

// serializeStringArray JSON配列をシリアライズするヘルパー関数
func serializeStringArray(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}
	data, err := json.Marshal(arr)
	if err != nil {
		return "[]"
	}
	return string(data)
}

