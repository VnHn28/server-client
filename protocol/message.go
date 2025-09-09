package protocol

import (
	"encoding/json"
	"time"
)

// Message Structs

type AuthMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TimeMessage struct {
	Timestamp time.Time `json:"timestamp"`
}

type AckMessage struct {
	Status string `json:"status"`
}

// Serialization Helpers

func Encode(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
