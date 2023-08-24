package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (session Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(session)
}

func (session *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, session)
}
