package events

import "encoding/json"

// Marshal event to json string
func MarshalEvent(e Event) (string, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

// Unmarshal event to given struct
func UnmarshalEvent(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}
