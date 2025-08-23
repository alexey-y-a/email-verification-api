package model

import "time"

type Verification struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Verified  bool      `json:"verified"`
}
