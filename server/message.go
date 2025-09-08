package server

import (
	"encoding/json"
	"time"
)

// AuthMessage is the initial message from client to server containing credentials.
type AuthMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TimeMessage is sent from client every 7 seconds after auth.
type TimeMessage struct {
	Timestamp time.Time `json:"timestamp"`
}

// AckMessage is sent by server to confirm receipt of TimeMessage.
type AckMessage struct {
	Status string `json:"status"`
}

// ---- Serialization Helpers ----

// ToJSON serializes a struct into JSON bytes.
func ToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// FromJSON deserializes JSON bytes into a struct.
func FromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
