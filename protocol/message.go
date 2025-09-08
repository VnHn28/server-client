package protocol

import (
	"time"
)

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
