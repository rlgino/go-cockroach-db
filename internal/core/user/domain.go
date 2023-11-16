package user

import "github.com/google/uuid"

type Data struct {
	ID       uuid.UUID `json:"id,omitempty"`
	User     string    `json:"user,omitempty"`
	Password string    `json:"password,omitempty"`
}
